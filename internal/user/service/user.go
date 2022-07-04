package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/devpies/saas-core/internal/user/model"
	"github.com/devpies/saas-core/pkg/msg"
	"github.com/devpies/saas-core/pkg/web"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"time"
)

type userRepository interface {
	RunTx(ctx context.Context, fn func(*sqlx.Tx) error) error
	CreateTx(ctx context.Context, tx *sqlx.Tx, nu model.NewUser, now time.Time) (model.User, error)
	CreateAdminTx(ctx context.Context, tx *sqlx.Tx, na model.NewAdminUser) error
	List(ctx context.Context) ([]model.User, error)
	RetrieveByEmail(ctx context.Context, email string) (model.User, error)
	RetrieveMe(ctx context.Context) (model.User, error)
	RemoveUserTx(ctx context.Context, tx *sqlx.Tx, uid string) error
}

type seatRepository interface {
	IncrementSeatsUsedTx(ctx context.Context, tx *sqlx.Tx) error
	DecrementSeatsUsedTx(ctx context.Context, tx *sqlx.Tx) error
	FindSeatsAvailable(ctx context.Context) (model.Seats, error)
	InsertSeatsEntryTx(ctx context.Context, tx *sqlx.Tx, maxSeats int) error
}

type cognitoClient interface {
	AdminCreateUser(
		ctx context.Context,
		params *cognitoidentityprovider.AdminCreateUserInput,
		optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminCreateUserOutput, error)
	AdminDeleteUser(
		ctx context.Context,
		params *cognitoidentityprovider.AdminDeleteUserInput,
		optFns ...func(*cognitoidentityprovider.Options),
	) (*cognitoidentityprovider.AdminDeleteUserOutput, error)
	AdminGetUser(
		ctx context.Context,
		params *cognitoidentityprovider.AdminGetUserInput,
		optFns ...func(*cognitoidentityprovider.Options),
	) (*cognitoidentityprovider.AdminGetUserOutput, error)
}

type connectionRepository interface {
	Insert(ctx context.Context, nc model.NewConnection) error
}

// UserService manages the user business operations.
type UserService struct {
	logger         *zap.Logger
	userPoolID     string
	userRepo       userRepository
	seatRepo       seatRepository
	cognitoClient  cognitoClient
	connectionRepo connectionRepository
}

const (
	MaximumSeatsBasic   = 3
	MaximumSeatsPremium = 25
)

// NewUserService returns a new user service.
func NewUserService(
	logger *zap.Logger,
	userPoolID string,
	userRepo userRepository,
	seatRepo seatRepository,
	cognitoClient cognitoClient,
	connectionRepo connectionRepository,
) *UserService {
	return &UserService{
		logger:         logger,
		userPoolID:     userPoolID,
		userRepo:       userRepo,
		seatRepo:       seatRepo,
		cognitoClient:  cognitoClient,
		connectionRepo: connectionRepo,
	}
}

// AddUser adds a new or existing user to the tenant's account and updates the number of seats.
func (us *UserService) AddUser(ctx context.Context, nu model.NewUser, now time.Time) (model.User, error) {
	getUserInput := &cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: aws.String(us.userPoolID),
		Username:   aws.String(nu.Email),
	}

	if _, err := us.cognitoClient.AdminGetUser(ctx, getUserInput); err != nil {
		var unf *types.UserNotFoundException

		if errors.As(err, &unf) {
			if err = us.createCognitoUser(ctx, nu); err != nil {
				return model.User{}, err
			}
			us.logger.Info("successfully added user")
		}
	}

	user, err := us.addUserToTenant(ctx, nu, now)
	if err != nil {
		return model.User{}, err
	}

	return user, err
}

func (us *UserService) addUserToTenant(ctx context.Context, nu model.NewUser, now time.Time) (model.User, error) {
	var (
		user model.User
		err  error
	)

	err = us.userRepo.RunTx(ctx, func(tx *sqlx.Tx) error {
		// Try creating a user for the tenant.
		user, err = us.userRepo.CreateTx(ctx, tx, nu, now)
		if err != nil {
			us.logger.Error("failed to add seat")
			return err
		}
		// Then increment seats used.
		if err = us.seatRepo.IncrementSeatsUsedTx(ctx, tx); err != nil {
			us.logger.Error("failed to increment seats used")
			return err
		}
		return nil
	})
	if err != nil {
		return model.User{}, err
	}

	if err = us.connectionRepo.Insert(ctx, model.NewConnection{UserID: user.ID, TenantID: ""}); err != nil {
		return user, err
	}

	return user, nil
}

func (us *UserService) createCognitoUser(ctx context.Context, nu model.NewUser) error {
	values, ok := web.FromContext(ctx)
	if !ok {
		return web.CtxErr()
	}
	_, err := us.cognitoClient.AdminCreateUser(ctx, &cognitoidentityprovider.AdminCreateUserInput{
		UserPoolId: aws.String(us.userPoolID),
		Username:   aws.String(nu.Email),
		UserAttributes: []types.AttributeType{
			{Name: aws.String("custom:tenant-id"), Value: aws.String(values.TenantID)},
			{Name: aws.String("custom:company-name"), Value: aws.String(nu.Company)},
			{Name: aws.String("custom:full-name"), Value: aws.String(fmt.Sprintf("%s %s", nu.FirstName, nu.LastName))},
			{Name: aws.String("email"), Value: aws.String(nu.Email)},
			{Name: aws.String("email_verified"), Value: aws.String("true")},
		},
	})
	return err
}

func (us *UserService) AddAdminUserFromEvent(ctx context.Context, message interface{}) error {
	m, err := msg.Bytes(message)
	if err != nil {
		return err
	}

	event, err := msg.UnmarshalTenantIdentityCreatedEvent(m)
	if err != nil {
		return err
	}

	na := newAdminUser(event.Data)

	// Create tenant admin and set max seats allowed.
	err = us.userRepo.RunTx(ctx, func(tx *sqlx.Tx) error {
		// When the subject is a tenant, the tenantID is already available in context.
		// In this case, the subject may be the SaaS Admin provider provisioning a tenant in the admin app.
		ctx = web.NewContext(ctx, &web.Values{TenantID: na.TenantID})

		if err = us.userRepo.CreateAdminTx(ctx, tx, na); err != nil {
			us.logger.Error("failed to create tenant admin")
			return err
		}
		// Determine max seats based on plan.
		if err = us.configureMaxSeats(ctx, tx, event.Data.Plan); err != nil {
			us.logger.Error("failed to insert seats entry")
			return err
		}
		return nil
	})
	return nil
}

func (us *UserService) configureMaxSeats(ctx context.Context, tx *sqlx.Tx, plan string) error {
	var maxSeats = MaximumSeatsBasic
	if plan == "premium" {
		maxSeats = MaximumSeatsPremium
	}
	// Add entry to seats table.
	return us.seatRepo.InsertSeatsEntryTx(ctx, tx, maxSeats)
}

func newAdminUser(data msg.TenantIdentityCreatedEventData) model.NewAdminUser {
	return model.NewAdminUser{
		UserID:        data.UserID,
		TenantID:      data.TenantID,
		Company:       data.Company,
		Email:         data.Email,
		FirstName:     data.FirstName,
		LastName:      data.LastName,
		EmailVerified: true,
		CreatedAt:     msg.ParseTime(data.CreatedAt),
	}
}

func (us *UserService) List(ctx context.Context) ([]model.User, error) {
	return us.userRepo.List(ctx)
}

func (us *UserService) RetrieveByEmail(ctx context.Context, email string) (model.User, error) {
	return us.userRepo.RetrieveByEmail(ctx, email)
}

func (us *UserService) RetrieveMe(ctx context.Context) (model.User, error) {
	return us.userRepo.RetrieveMe(ctx)
}

func (us *UserService) SeatsAvailable(ctx context.Context) (model.SeatsAvailableResult, error) {
	var res model.SeatsAvailableResult

	result, err := us.seatRepo.FindSeatsAvailable(ctx)
	if err != nil {
		return res, err
	}

	seatsAvailable := result.MaxSeats - result.SeatsUsed
	if seatsAvailable < 0 {
		us.logger.Error("seats available is less than 0")
		return res, web.NewShutdownError("unexpected seats available")
	}

	return model.SeatsAvailableResult{
		MaxSeats:       result.MaxSeats,
		SeatsAvailable: seatsAvailable,
	}, nil
}

func (us *UserService) RemoveUser(ctx context.Context, uid, email string) error {
	// Remove user and decrement the seats used counter.
	if err := us.userRepo.RunTx(ctx, func(tx *sqlx.Tx) error {
		// Remove user.
		if err := us.userRepo.RemoveUserTx(ctx, tx, uid); err != nil {
			us.logger.Error("failed to remove user")
			return err
		}
		// Decrement seats used.
		if err := us.seatRepo.DecrementSeatsUsedTx(ctx, tx); err != nil {
			us.logger.Error("failed to decrement seats used")
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	if _, err := us.cognitoClient.AdminDeleteUser(ctx, &cognitoidentityprovider.AdminDeleteUserInput{
		Username:   aws.String(email),
		UserPoolId: aws.String(us.userPoolID),
	}); err != nil {
		us.logger.Error("failed to remove cognito user", zap.Error(err))
		return err
	}
	us.logger.Info("successfully removed cognito user")

	return nil
}

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
	"net/http"
	"time"
)

type userRepository interface {
	RunTx(ctx context.Context, fn func(*sqlx.Tx) error) error
	AddUserTx(ctx context.Context, tx *sqlx.Tx, userID string, now time.Time) error
	CreateUserProfile(ctx context.Context, nu model.NewUser, userID string, now time.Time) (model.User, error)
	CreateAdminUser(ctx context.Context, na model.NewAdminUser) error
	List(ctx context.Context) ([]model.User, error)
	RetrieveIDByEmail(ctx context.Context, email string) (string, error)
	RetrieveByEmail(ctx context.Context, email string) (model.User, error)
	RetrieveMe(ctx context.Context) (model.User, error)
	DetachUserTx(ctx context.Context, tx *sqlx.Tx, uid string) error
}

type seatRepository interface {
	IncrementSeatsUsedTx(ctx context.Context, tx *sqlx.Tx) error
	DecrementSeatsUsedTx(ctx context.Context, tx *sqlx.Tx) error
	FindSeatsAvailable(ctx context.Context) (model.Seats, error)
	InsertSeatsEntryTx(ctx context.Context, tx *sqlx.Tx, maxSeats model.MaximumSeatsType) error
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
	Delete(ctx context.Context, userID string) error
}

// UserService manages the user business operations.
type UserService struct {
	logger           *zap.Logger
	userRepo         userRepository
	seatRepo         seatRepository
	cognitoClient    cognitoClient
	connectionRepo   connectionRepository
	sharedUserPoolID string
}

// NewUserService returns a new user service.
func NewUserService(
	logger *zap.Logger,
	userRepo userRepository,
	seatRepo seatRepository,
	cognitoClient cognitoClient,
	connectionRepo connectionRepository,
	sharedUserPoolID string,
) *UserService {
	return &UserService{
		logger:           logger,
		userRepo:         userRepo,
		seatRepo:         seatRepo,
		cognitoClient:    cognitoClient,
		connectionRepo:   connectionRepo,
		sharedUserPoolID: sharedUserPoolID,
	}
}

// AddUser adds a new or existing user to the tenant's account and updates the number of seats.
func (us *UserService) AddUser(ctx context.Context, nu model.NewUser, now time.Time) error {
	var (
		userID string
		err    error
	)

	values, ok := web.FromContext(ctx)
	if !ok {
		return web.CtxErr()
	}

	// Verify user doesn't belong to the tenant already.
	_, err = us.userRepo.RetrieveByEmail(ctx, nu.Email)
	if err == nil {
		us.logger.Info("user already connected to tenant")
		return web.NewRequestError(fmt.Errorf("user already added"), http.StatusBadRequest)
	}

	// Determine if a cognito identity already exists.
	output, err := us.cognitoClient.AdminGetUser(ctx, &cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: aws.String(us.sharedUserPoolID),
		Username:   aws.String(nu.Email),
	})
	// If identity exists, attach existing user profile to tenant.
	if err == nil {
		userID = getUserIDFromAttributes(output.UserAttributes)
		err = us.addUserToTenant(ctx, userID, values.TenantID, now)
		if err != nil {
			return err
		}
		return nil
	}

	// Otherwise, create identity when user not found, then add user to tenant.
	var unf *types.UserNotFoundException

	if errors.As(err, &unf) {
		var result *cognitoidentityprovider.AdminCreateUserOutput

		result, err = us.createCognitoIdentity(ctx, nu)
		if err != nil {
			us.logger.Error("error creating cognito identity")
			return err
		}
		userID = getUserIDFromAttributes(result.User.Attributes)
		_, err = us.userRepo.CreateUserProfile(ctx, nu, userID, now)
		if err != nil {
			us.logger.Error("error creating user profile")
			return err
		}
		err = us.addUserToTenant(ctx, userID, values.TenantID, now)
		if err != nil {
			return err
		}
	}

	return nil
}

func getUserIDFromAttributes(attributes []types.AttributeType) string {
	for _, v := range attributes {
		if v.Name != nil && *v.Name == "sub" {
			return *v.Value
		}
	}
	return ""
}

func (us *UserService) addUserToTenant(ctx context.Context, userID, tenantID string, now time.Time) error {
	var err error

	err = us.userRepo.RunTx(ctx, func(tx *sqlx.Tx) error {
		// Add user to tenant.
		err = us.userRepo.AddUserTx(ctx, tx, userID, now)
		if err != nil {
			us.logger.Error("error adding user to tenant")
			return err
		}
		// Then increment seats used.
		if err = us.seatRepo.IncrementSeatsUsedTx(ctx, tx); err != nil {
			us.logger.Error("error incrementing seats used")
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	connection := model.NewConnection{
		UserID:   userID,
		TenantID: tenantID,
	}

	if err = us.connectionRepo.Insert(ctx, connection); err != nil {
		us.logger.Error("error creating connection")
		return err
	}

	return nil
}

func (us *UserService) createCognitoIdentity(ctx context.Context, nu model.NewUser) (*cognitoidentityprovider.AdminCreateUserOutput, error) {
	values, ok := web.FromContext(ctx)
	if !ok {
		return nil, web.CtxErr()
	}
	return us.cognitoClient.AdminCreateUser(ctx, &cognitoidentityprovider.AdminCreateUserInput{
		UserPoolId: aws.String(us.sharedUserPoolID),
		Username:   aws.String(nu.Email),
		UserAttributes: []types.AttributeType{
			{Name: aws.String("custom:tenant-id"), Value: aws.String(values.TenantID)},
			{Name: aws.String("custom:company-name"), Value: aws.String(nu.Company)},
			{Name: aws.String("custom:full-name"), Value: aws.String(fmt.Sprintf("%s %s", nu.FirstName, nu.LastName))},
			{Name: aws.String("email"), Value: aws.String(nu.Email)},
			{Name: aws.String("email_verified"), Value: aws.String("true")},
		},
	})
}

// AddAdminUserFromEvent creates the tenant admin and sets max seats allowed.
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

	err = us.userRepo.RunTx(ctx, func(tx *sqlx.Tx) error {
		// The requester may be a SaaS Admin. Use the tenantID from the event instead of context.
		ctx = web.NewContext(ctx, &web.Values{TenantID: na.TenantID})

		if err = us.userRepo.CreateAdminUser(ctx, na); err != nil {
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
	var maxSeats = model.MaximumSeatsBasic
	if plan == "premium" {
		maxSeats = model.MaximumSeatsPremium
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

// List returns all users associated to the tenant account.
func (us *UserService) List(ctx context.Context) ([]model.User, error) {
	return us.userRepo.List(ctx)
}

// RetrieveByEmail retrieves a user by email.
func (us *UserService) RetrieveByEmail(ctx context.Context, email string) (model.User, error) {
	return us.userRepo.RetrieveByEmail(ctx, email)
}

// RetrieveMe returns the requester user details.
func (us *UserService) RetrieveMe(ctx context.Context) (model.User, error) {
	return us.userRepo.RetrieveMe(ctx)
}

// SeatsAvailable returns the number of remaining seats available.
func (us *UserService) SeatsAvailable(ctx context.Context) (model.SeatsAvailableResult, error) {
	var res model.SeatsAvailableResult

	result, err := us.seatRepo.FindSeatsAvailable(ctx)
	if err != nil {
		return res, err
	}

	seatsAvailable := int(result.MaxSeats) - result.SeatsUsed
	if seatsAvailable < 0 {
		us.logger.Error("seats available is less than 0")
		return res, web.NewShutdownError("unexpected seats available")
	}

	return model.SeatsAvailableResult{
		MaxSeats:       result.MaxSeats,
		SeatsAvailable: seatsAvailable,
	}, nil
}

// RemoveUser removes a user from a tenant account and updates the available seats.
func (us *UserService) RemoveUser(ctx context.Context, uid string) error {
	// Remove user and decrement the seats used counter.
	if err := us.userRepo.RunTx(ctx, func(tx *sqlx.Tx) error {
		// Remove user.
		if err := us.userRepo.DetachUserTx(ctx, tx, uid); err != nil {
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

	return us.connectionRepo.Delete(ctx, uid)
}

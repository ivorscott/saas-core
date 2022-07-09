package service

import (
	"context"
	"fmt"

	"github.com/alexedwards/scs/v2"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"go.uber.org/zap"
)

// AuthService is responsible for managing authentication with Cognito.
type AuthService struct {
	logger          *zap.Logger
	region          string
	adminClientID   string
	adminUserPoolID string
	cognitoClient   cognitoClient
	session         *scs.SessionManager
}

// NewAuthService creates a new instance of AuthService.
func NewAuthService(logger *zap.Logger, region, adminClientID, adminUserPoolID string, cognitoClient cognitoClient, session *scs.SessionManager) *AuthService {
	return &AuthService{
		logger:          logger,
		region:          region,
		adminClientID:   adminClientID,
		adminUserPoolID: adminUserPoolID,
		cognitoClient:   cognitoClient,
		session:         session,
	}
}

// Authenticate initiates the server-side auth flow.
func (as *AuthService) Authenticate(ctx context.Context, email, password string) (*cognitoidentityprovider.AdminInitiateAuthOutput, error) {
	signInInput := &cognitoidentityprovider.AdminInitiateAuthInput{
		AuthFlow:       "ADMIN_USER_PASSWORD_AUTH",
		ClientId:       aws.String(as.adminClientID),
		UserPoolId:     aws.String(as.adminUserPoolID),
		AuthParameters: map[string]string{"USERNAME": email, "PASSWORD": password},
	}

	return as.cognitoClient.AdminInitiateAuth(ctx, signInInput)
}

// RespondToNewPasswordRequiredChallenge completes the server-side auth flow for freshly onboarded users.
func (as *AuthService) RespondToNewPasswordRequiredChallenge(ctx context.Context, email, password string, session string) (*cognitoidentityprovider.AdminRespondToAuthChallengeOutput, error) {
	params := &cognitoidentityprovider.AdminRespondToAuthChallengeInput{
		ChallengeName: "NEW_PASSWORD_REQUIRED",
		ClientId:      aws.String(as.adminClientID),
		UserPoolId:    aws.String(as.adminUserPoolID),
		ChallengeResponses: map[string]string{
			"USERNAME":     email,
			"NEW_PASSWORD": password,
		},
		Session: aws.String(session),
	}
	return as.cognitoClient.AdminRespondToAuthChallenge(ctx, params)
}

// CreateUserSession parses the idToken and saves the subject.
func (as *AuthService) CreateUserSession(ctx context.Context, idToken []byte) error {
	pubKeyURL := "https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json"
	formattedURL := fmt.Sprintf(pubKeyURL, as.region, as.adminUserPoolID)

	keySet, err := jwk.Fetch(ctx, formattedURL)
	if err != nil {
		as.logger.Error("error fetching token", zap.Error(err))
		return err
	}

	tok, err := jwt.Parse(idToken, jwt.WithKeySet(keySet))
	if err != nil {
		as.logger.Error("error decoding token", zap.Error(err))
		return err
	}

	// Retrieve values.
	sub := tok.Subject()
	email, ok := tok.Get("email")
	if !ok {
		return fmt.Errorf("email is not available")
	}

	// Store session.
	as.session.Put(ctx, "UserID", sub)
	as.session.Put(ctx, "Email", email.(string))

	return nil
}

// CreatePasswordChallengeSession creates a session for the active password challenge.
// This is used to deny access to the change password form. Only users with an active password challenge can view it.
func (as *AuthService) CreatePasswordChallengeSession(ctx context.Context) {
	as.session.Put(ctx, "PasswordChallenge", true)
}

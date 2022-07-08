package web

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"go.uber.org/zap"
)

var ErrInvalidAuthorizationHeader = errors.New("missing or invalid authorization header")

func Authenticate(log *zap.Logger, r *http.Request, region string, userPoolID string) (*http.Request, error) {
	authHeader := r.Header.Get("Authorization")
	token, sub, tenantID, tenantMap, err := verifyToken(r.Context(), log, authHeader, region, userPoolID)
	if err != nil {
		return nil, NewRequestError(err, http.StatusUnauthorized)
	}
	r = addContextMetadata(r, token, sub, tenantID, tenantMap)
	return r, nil
}

func getToken(authHeader string) (string, error) {
	if len(authHeader) > 7 && strings.ToLower(authHeader[0:6]) == "bearer" {
		return authHeader[7:], nil
	}
	return "", ErrInvalidAuthorizationHeader
}

// TenantConnectionMap represents a valid tenant connection mapping.
type TenantConnectionMap map[string]struct {
	TenantID    string `json:"id"`
	CompanyName string `json:"companyName"`
	Plan        string `json:"plan"`
	Path        string `json:"path"`
}

func verifyToken(
	ctx context.Context,
	logger *zap.Logger,
	authHeader string,
	region string,
	userPoolID string,
) (token string, sub string, tenantID string, tenantMap TenantConnectionMap, err error) {
	token, err = getToken(authHeader)
	if err != nil {
		return
	}

	pubKeyURL := "https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json"
	formattedURL := fmt.Sprintf(pubKeyURL, region, userPoolID)

	keySet, err := jwk.Fetch(ctx, formattedURL)
	if err != nil {
		logger.Error("error fetching token", zap.String("pubKeyURL", formattedURL), zap.Error(err))
		return
	}

	parsedToken, err := jwt.Parse(
		[]byte(token),
		jwt.WithKeySet(keySet),
		jwt.WithValidate(true),
	)
	if err != nil {
		logger.Error("error decoding token", zap.Error(err))
		return
	}
	sub = parsedToken.Subject()

	val, ok := parsedToken.Get("custom:tenant-id")
	if ok {
		tenantID = val.(string)
	}

	val, ok = parsedToken.Get("custom:tenant-connections")
	if ok {
		err = json.Unmarshal([]byte(val.(string)), &tenantMap)
		if err != nil {
			return
		}
	}

	return
}

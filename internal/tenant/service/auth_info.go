package service

import (
	"context"
	"github.com/devpies/saas-core/internal/tenant/model"
	"go.uber.org/zap"
	"regexp"
)

// AuthInfoService retrieves tenant auth information.
type AuthInfoService struct {
	logger       *zap.Logger
	region       string
	authInfoRepo authInfoRepository
}

type authInfoRepository interface {
	SelectAuthInfo(ctx context.Context, path string) (model.AuthInfo, error)
}

// NewAuthInfoService returns a new AuthInfoService.
func NewAuthInfoService(logger *zap.Logger, authInfoRepo authInfoRepository, region string) *AuthInfoService {
	return &AuthInfoService{
		logger:       logger,
		region:       region,
		authInfoRepo: authInfoRepo,
	}
}

// GetInfo gets the tenant authentication information.
func (ts *AuthInfoService) GetInfo(ctx context.Context, referer string) (model.AuthInfoAndRegion, error) {
	var info model.AuthInfoAndRegion
	path := getPath(referer)
	result, err := ts.authInfoRepo.SelectAuthInfo(ctx, path)
	if err != nil {
		return info, err
	}
	info = model.AuthInfoAndRegion{
		ProjectRegion:    ts.region,
		CognitoRegion:    ts.region,
		UserPoolID:       result.UserPoolID,
		UserPoolClientID: result.UserPoolClientID,
	}
	return info, nil
}

// getPath parses the request URI and retrieves the base path. The base path is either "app" or the shortened tenant name.
func getPath(referer string) string {
	var url = "localhost"
	if referer != "" {
		url = referer
	}
	var re = regexp.MustCompile(`(?i)\/?(\w+)(?:\/index\.html)?`)
	return re.FindAllString(url, -1)[1]
}

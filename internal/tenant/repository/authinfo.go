package repository

import (
	"context"

	"github.com/devpies/saas-core/internal/tenant/model"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// AuthInfoRepository manages data access to tenant authentication information.
type AuthInfoRepository struct {
	client *dynamodb.Client
}

// NewAuthInfoRepository returns a new AuthInfoRepository.
func NewAuthInfoRepository(client *dynamodb.Client) *AuthInfoRepository {
	return &AuthInfoRepository{
		client: client,
	}
}

// SelectAuthInfo retrieves authentication information for a specific tenant.
func (ar *AuthInfoRepository) SelectAuthInfo(ctx context.Context, path string) (model.AuthInfo, error) {
	var authInfo model.AuthInfo
	return authInfo, nil
}

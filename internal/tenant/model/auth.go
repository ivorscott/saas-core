package model

import (
	"github.com/go-playground/validator/v10"
)

var authInfoValidator *validator.Validate

func init() {
	v := NewValidator()
	authInfoValidator = v
}

// AuthInfo represents tenant authentication information.
type AuthInfo struct {
	ProjectRegion    string `json:"projectRegion"`
	CognitoRegion    string `json:"cognitoRegion"`
	UserPoolID       string `json:"userPoolId"`
	UserPoolClientID string `json:"userPoolClientId"`
}

// Validate validates the AuthInfo.
func (a *AuthInfo) Validate() error {
	return authInfoValidator.Struct(a)
}

package model_test

import (
	"testing"

	"github.com/devpies/saas-core/internal/admin/model"

	"github.com/stretchr/testify/assert"
)

func TestAuthCredentials_Validate(t *testing.T) {
	tests := []struct {
		name     string
		modifier func(ac *model.AuthCredentials)
		err      string
	}{
		{
			name:     "valid",
			modifier: func(ac *model.AuthCredentials) {},
			err:      "",
		},
		{
			name: "invalid: missing email",
			modifier: func(ac *model.AuthCredentials) {
				ac.Email = ""
			},
			err: "failed on the 'required' tag",
		},
		{
			name: "invalid: missing password",
			modifier: func(ac *model.AuthCredentials) {
				ac.Password = ""
			},
			err: "failed on the 'required' tag",
		},
		{
			name: "invalid: not an email",
			modifier: func(ac *model.AuthCredentials) {
				ac.Email = "notAnEmail"
			},
			err: "failed on the 'email' tag",
		},
		{
			name: "invalid: password character length",
			modifier: func(ac *model.AuthCredentials) {
				ac.Password = "pass"
			},
			err: "failed on the 'min' tag",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ac := model.AuthCredentials{
				Email:    "test@email.com",
				Password: "P@ssw0rd",
			}

			tc.modifier(&ac)

			err := ac.Validate()
			if tc.err != "" {
				assert.Regexp(t, tc.err, err.Error())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

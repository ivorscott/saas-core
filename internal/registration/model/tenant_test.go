package model_test

import (
	"testing"

	"github.com/devpies/saas-core/internal/registration/model"

	"github.com/stretchr/testify/assert"
)

func TestNewTenant_Validate(t *testing.T) {
	tests := []struct {
		name     string
		modifier func(nt *model.NewTenant)
		err      string
	}{
		{
			name:     "valid",
			modifier: func(nt *model.NewTenant) {},
			err:      "",
		},
		{
			name: "valid premium tenant",
			modifier: func(nt *model.NewTenant) {
				nt.Plan = "premium"
			},
			err: "",
		},
		{
			name: "invalid id",
			modifier: func(nt *model.NewTenant) {
				nt.ID = ""
			},
			err: "failed on the 'required' tag",
		},
		{
			name: "invalid first name",
			modifier: func(nt *model.NewTenant) {
				nt.FirstName = ""
			},
			err: "failed on the 'required' tag",
		},
		{
			name: "invalid last name",
			modifier: func(nt *model.NewTenant) {
				nt.LastName = ""
			},
			err: "failed on the 'required' tag",
		},
		{
			name: "invalid company",
			modifier: func(nt *model.NewTenant) {
				nt.Company = ""
			},
			err: "failed on the 'required' tag",
		},
		{
			name: "invalid email",
			modifier: func(nt *model.NewTenant) {
				nt.Email = ""
			},
			err: "failed on the 'required' tag",
		},
		{
			name: "invalid plan",
			modifier: func(nt *model.NewTenant) {
				nt.Plan = ""
			},
			err: "failed on the 'required' tag",
		},
		{
			name: "unknown plan",
			modifier: func(nt *model.NewTenant) {
				nt.Plan = "enterprise"
			},
			err: "failed on the 'oneof' tag",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			nt := model.NewTenant{
				ID:        "ca5b783c-3c3a-401f-b135-b42a2dfec5d5",
				FirstName: "Test",
				LastName:  "User",
				Company:   "Test Company",
				Email:     "test@email.com",
				Plan:      "basic",
			}

			tc.modifier(&nt)

			err := nt.Validate()
			if tc.err != "" {
				if err == nil {
					t.Errorf("expected: %s, got nil", tc.err)
					return
				}
				assert.Regexp(t, tc.err, err.Error())
			} else {
				if err != nil {
					t.Errorf("expected: nil, got: %s", err.Error())
				}
			}
		})
	}
}

package model_test

import (
	"github.com/devpies/saas-core/internal/admin/model"
	"github.com/stretchr/testify/assert"
	"testing"
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
			name: "valid (premium tenant)",
			modifier: func(nt *model.NewTenant) {
				nt.Plan = "premium"
			},
			err: "",
		},
		{
			name: "invalid full name",
			modifier: func(nt *model.NewTenant) {
				nt.FullName = ""
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
				FullName: "Test User",
				Company:  "Test Company",
				Email:    "test@email.com",
				Plan:     "basic",
			}

			tc.modifier(&nt)

			err := nt.Validate()
			if tc.err != "" {
				assert.Regexp(t, tc.err, err.Error())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

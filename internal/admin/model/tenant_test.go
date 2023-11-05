package model_test

import (
	"testing"

	"github.com/devpies/saas-core/internal/admin/model"

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
				nt.CompanyName = ""
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
				ID:          "84c59228-52ff-44a2-8163-a429de9fe77b",
				FirstName:   "Test",
				LastName:    "User",
				CompanyName: "Test Company",
				Email:       "test@email.com",
				Plan:        "basic",
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

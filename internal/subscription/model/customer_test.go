package model_test

import (
	"testing"
	"time"

	"github.com/devpies/saas-core/internal/subscription/model"

	"github.com/stretchr/testify/assert"
)

func TestNewCustomer_Validate(t *testing.T) {
	tests := []struct {
		name     string
		modifier func(nc *model.NewCustomer)
		err      string
	}{
		{
			name:     "valid",
			modifier: func(nc *model.NewCustomer) {},
			err:      "",
		},
		{
			name: "invalid first name",
			modifier: func(nc *model.NewCustomer) {
				nc.FirstName = ""
			},
			err: "failed on the 'required' tag",
		},
		{
			name: "invalid first name too long",
			modifier: func(nc *model.NewCustomer) {
				var text string
				for i := 0; i < 256; i++ {
					text += "x"
				}
				nc.FirstName = text
			},
			err: "failed on the 'max' tag",
		},
		{
			name: "invalid last name",
			modifier: func(nc *model.NewCustomer) {
				nc.LastName = ""
			},
			err: "failed on the 'required' tag",
		},
		{
			name: "invalid last name too long",
			modifier: func(nc *model.NewCustomer) {
				var text string
				for i := 0; i < 256; i++ {
					text += "x"
				}
				nc.LastName = text
			},
			err: "failed on the 'max' tag",
		},
		{
			name: "invalid payment method",
			modifier: func(nc *model.NewCustomer) {
				nc.PaymentMethodID = ""
			},
			err: "failed on the 'required' tag",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			nc := model.NewCustomer{
				FirstName:       "Julian",
				LastName:        "Smith",
				Email:           "julian@example.com",
				ID:              "cus_OxwmOOTrGdKGcH",
				PaymentMethodID: "pm_1KzRIfIbOZLMWfd3M0YVrVfc",
			}

			tc.modifier(&nc)

			err := nc.Validate()
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

func TestUpdateCustomer_Validate(t *testing.T) {
	tests := []struct {
		name     string
		modifier func(uc *model.UpdateCustomer)
		err      string
	}{
		{
			name:     "valid",
			modifier: func(uc *model.UpdateCustomer) {},
			err:      "",
		},
		{
			name: "invalid first name",
			modifier: func(uc *model.UpdateCustomer) {
				name := ""
				uc.FirstName = &name
			},
			err: "failed on the 'min' tag",
		},
		{
			name: "invalid last name",
			modifier: func(uc *model.UpdateCustomer) {
				name := ""
				uc.LastName = &name
			},
			err: "failed on the 'min' tag",
		},
		{
			name: "invalid email",
			modifier: func(uc *model.UpdateCustomer) {
				email := "not-email"
				uc.Email = &email
			},
			err: "failed on the 'email' tag",
		},
		{
			name: "invalid update time",
			modifier: func(uc *model.UpdateCustomer) {
				uc.UpdatedAt = time.Time{}
			},
			err: "failed on the 'required' tag",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			uc := model.UpdateCustomer{
				UpdatedAt: time.Now(),
			}

			tc.modifier(&uc)

			err := uc.Validate()
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

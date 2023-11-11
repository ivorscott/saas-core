package model_test

import (
	"testing"
	"time"

	"github.com/devpies/saas-core/internal/subscription/model"

	"github.com/stretchr/testify/assert"
)

func TestNewTransaction_Validate(t *testing.T) {
	tests := []struct {
		name     string
		modifier func(nt *model.NewTransaction)
		err      string
	}{
		{
			name:     "valid",
			modifier: func(nt *model.NewTransaction) {},
			err:      "",
		},
		{
			name: "invalid amount",
			modifier: func(nt *model.NewTransaction) {
				nt.Amount = 0
			},
			err: "failed on the 'required' tag",
		},
		{
			name: "invalid currency",
			modifier: func(nt *model.NewTransaction) {
				nt.Currency = ""
			},
			err: "failed on the 'required' tag",
		},
		{
			name: "invalid last four digits",
			modifier: func(nt *model.NewTransaction) {
				nt.LastFour = ""
			},
			err: "failed on the 'required' tag",
		},
		{
			name: "invalid last four digits length",
			modifier: func(nt *model.NewTransaction) {
				nt.LastFour = "00001"
			},
			err: "failed on the 'len' tag",
		},
		{
			name: "invalid status id",
			modifier: func(nt *model.NewTransaction) {
				nt.StatusID = 5
			},
			err: "failed on the 'oneof' tag",
		},
		{
			name: "invalid expiration month",
			modifier: func(nt *model.NewTransaction) {
				nt.ExpirationMonth = 0
			},
			err: "failed on the 'gte' tag",
		},
		{
			name: "non existent expiration month",
			modifier: func(nt *model.NewTransaction) {
				nt.ExpirationMonth = 13
			},
			err: "failed on the 'lte' tag",
		},
		{
			name: "non existent expiration year",
			modifier: func(nt *model.NewTransaction) {
				nt.ExpirationYear = 1957
			},
			err: "failed on the 'min' tag",
		},
		{
			name: "invalid stripe subscription id",
			modifier: func(nt *model.NewTransaction) {
				nt.SubscriptionID = ""
			},
			err: "failed on the 'required' tag",
		},
		{
			name: "invalid payment method",
			modifier: func(nt *model.NewTransaction) {
				nt.PaymentMethod = ""
			},
			err: "failed on the 'required' tag",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			nt := model.NewTransaction{
				ID:              "txn_3OArDkIbOZLMWfd30wyBVTpt",
				Amount:          1000,
				Currency:        "eur",
				LastFour:        "0001",
				BankReturnCode:  "",
				StatusID:        2,
				ExpirationMonth: 2,
				ExpirationYear:  2025,
				SubscriptionID:  "sub_1KzXfrIbOZLMWfd3glCgX807",
				PaymentIntent:   "",
				PaymentMethod:   "pm_1KzZYPIbOZLMWfd3GeQlt32W",
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

func TestUpdateTransaction_Validate(t *testing.T) {
	tests := []struct {
		name     string
		modifier func(ut *model.UpdateTransaction)
		err      string
	}{
		{
			name:     "valid",
			modifier: func(ut *model.UpdateTransaction) {},
			err:      "",
		},
		{
			name: "invalid transaction status id",
			modifier: func(ut *model.UpdateTransaction) {
				var id model.TransactionStatusType = 5
				ut.StatusID = id
			},
			err: "failed on the 'oneof' tag",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ut := model.UpdateTransaction{
				StatusID:  2,
				UpdatedAt: time.Now(),
			}

			tc.modifier(&ut)

			err := ut.Validate()
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

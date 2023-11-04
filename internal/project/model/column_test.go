package model_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/devpies/saas-core/internal/project/model"
	"github.com/stretchr/testify/assert"
)

func TestNewColumn_Validate(t *testing.T) {
	tests := []struct {
		name     string
		modifier func(nc *model.NewColumn)
		err      string
	}{
		{
			name:     "valid",
			modifier: func(nc *model.NewColumn) {},
			err:      "",
		},
		{
			name: "missing title",
			modifier: func(nc *model.NewColumn) {
				nc.Title = ""
			},
			err: "failed on the 'required' tag",
		},
		{
			name: "missing column name",
			modifier: func(nc *model.NewColumn) {
				nc.ColumnName = ""
			},
			err: "failed on the 'required' tag",
		},
		{
			name: "missing project id",
			modifier: func(nc *model.NewColumn) {
				nc.ProjectID = ""
			},
			err: "failed on the 'required' tag",
		},
		{
			name: "title is too long",
			modifier: func(nc *model.NewColumn) {
				var text string
				for i := 0; i < 25; i++ {
					text += "x"
				}
				nc.Title = text
			},
			err: "failed on the 'max' tag",
		},
		{
			name: "column name is too long",
			modifier: func(nc *model.NewColumn) {
				var text string
				for i := 0; i < 11; i++ {
					text += "x"
				}
				nc.ColumnName = text
			},
			err: "failed on the 'max' tag",
		},
		{
			name: "project id is too long",
			modifier: func(nc *model.NewColumn) {
				var text string
				for i := 0; i < 37; i++ {
					text += "x"
				}
				nc.ProjectID = text
			},
			err: "failed on the 'max' tag",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			nc := model.NewColumn{
				Title:      "title",
				ColumnName: "column-x",
				ProjectID:  "project-id",
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

func TestUpdateColumn_Validate(t *testing.T) {
	tests := []struct {
		name     string
		modifier func(uc *model.UpdateColumn)
		err      string
	}{
		{
			name:     "valid",
			modifier: func(uc *model.UpdateColumn) {},
			err:      "",
		},
		{
			name: "title is too long",
			modifier: func(uc *model.UpdateColumn) {
				var text string
				for i := 0; i < 25; i++ {
					text += "x"
				}
				uc.Title = aws.String(text)
			},
			err: "failed on the 'max' tag",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			uc := model.UpdateColumn{}

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

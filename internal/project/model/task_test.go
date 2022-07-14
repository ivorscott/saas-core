package model_test

import (
	"testing"

	"github.com/devpies/saas-core/internal/project/model"

	"github.com/stretchr/testify/assert"
)

func TestNewTask_Validate(t *testing.T) {
	tests := []struct {
		name     string
		modifier func(nt *model.NewTask)
		err      string
	}{
		{
			name:     "valid",
			modifier: func(nt *model.NewTask) {},
			err:      "",
		},
		{
			name: "valid",
			modifier: func(nt *model.NewTask) {
				nt.Title = "This is an extremely long name for a title but should be accepted as valid."
			},
			err: "",
		},
		{
			name: "too long",
			modifier: func(nt *model.NewTask) {
				nt.Title = "This is an extremely long name for a title and should not be accepted as is!" // over by 1 char
			},
			err: "failed on the 'max' tag",
		},
		{
			name: "title empty",
			modifier: func(nt *model.NewTask) {
				nt.Title = ""
			},
			err: "failed on the 'required' tag",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			nt := model.NewTask{
				Title: "Test",
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

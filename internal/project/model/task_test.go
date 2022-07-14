package model_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
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
			name: "title is too long",
			modifier: func(nt *model.NewTask) {
				var text string
				for i := 0; i < 76; i++ {
					text += "x"
				}
				nt.Title = text
			},
			err: "failed on the 'max' tag",
		},
		{
			name: "missing title",
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

func TestUpdateTask_Validate(t *testing.T) {
	tests := []struct {
		name     string
		modifier func(ut *model.UpdateTask)
		err      string
	}{
		{
			name:     "valid",
			modifier: func(ut *model.UpdateTask) {},
			err:      "",
		},
		{
			name: "title too long",
			modifier: func(ut *model.UpdateTask) {
				var text string
				for i := 0; i < 76; i++ {
					text += "x"
				}
				ut.Title = aws.String(text)
			},
			err: "failed on the 'max' tag",
		},
		{
			name: "content too long",
			modifier: func(ut *model.UpdateTask) {
				var text string
				for i := 0; i < 1001; i++ {
					text += "x"
				}
				ut.Content = aws.String(text)
			},
			err: "failed on the 'max' tag",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ut := model.UpdateTask{}

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

func TestMoveTask_Validate(t *testing.T) {
	tests := []struct {
		name     string
		modifier func(mt *model.MoveTask)
		err      string
	}{
		{
			name:     "valid",
			modifier: func(mt *model.MoveTask) {},
			err:      "",
		},
		{
			name: "missing to",
			modifier: func(mt *model.MoveTask) {
				mt.To = ""
			},
			err: "failed on the 'required' tag",
		},
		{
			name: "missing from",
			modifier: func(mt *model.MoveTask) {
				mt.From = ""
			},
			err: "failed on the 'required' tag",
		},
		{
			name: "missing taskIds",
			modifier: func(mt *model.MoveTask) {
				mt.TaskIds = nil
			},
			err: "failed on the 'required' tag",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mt := model.MoveTask{
				To:      "to-column",
				From:    "from-column",
				TaskIds: []string{"task-id"},
			}

			tc.modifier(&mt)

			err := mt.Validate()
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

package model_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/devpies/saas-core/internal/project/model"
	"github.com/stretchr/testify/assert"
)

func TestNewComment_Validate(t *testing.T) {
	tests := []struct {
		name     string
		modifier func(nc *model.NewComment)
		err      string
	}{
		{
			name:     "valid",
			modifier: func(nc *model.NewComment) {},
			err:      "",
		},
		{
			name: "missing content",
			modifier: func(nc *model.NewComment) {
				nc.Content = ""
			},
			err: "failed on the 'required' tag",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			nc := model.NewComment{
				Content: "comment",
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

func TestUpdateComment_Validate(t *testing.T) {
	tests := []struct {
		name     string
		modifier func(uc *model.UpdateComment)
		err      string
	}{
		{
			name:     "valid",
			modifier: func(uc *model.UpdateComment) {},
			err:      "",
		},
		{
			name: "content is too long",
			modifier: func(uc *model.UpdateComment) {
				var text string
				for i := 0; i < 501; i++ {
					text += "x"
				}
				uc.Content = aws.String(text)
			},
			err: "failed on the 'max' tag",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			uc := model.UpdateComment{}

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

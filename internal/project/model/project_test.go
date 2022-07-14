package model_test

import (
	"testing"

	"github.com/devpies/saas-core/internal/project/model"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"
)

func TestNewProject_Validate(t *testing.T) {
	tests := []struct {
		name     string
		modifier func(np *model.NewProject)
		err      string
	}{
		{
			name:     "valid",
			modifier: func(np *model.NewProject) {},
			err:      "",
		},
		{
			name: "too long",
			modifier: func(np *model.NewProject) {
				np.Name = "Extremely Long Name For A Project"
			},
			err: "failed on the 'max' tag",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			np := model.NewProject{
				Name: "Test",
			}

			tc.modifier(&np)

			err := np.Validate()
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

func TestUpdateProject_Validate(t *testing.T) {
	tests := []struct {
		name     string
		modifier func(up *model.UpdateProject)
		err      string
	}{
		{
			name:     "valid",
			modifier: func(up *model.UpdateProject) {},
			err:      "",
		},
		{
			name: "valid",
			modifier: func(up *model.UpdateProject) {
				up.Description = aws.String("This text represents a project description that is really long but valid")
			},
			err: "",
		},
		{
			name: "too long",
			modifier: func(up *model.UpdateProject) {
				up.Name = aws.String("extremely long project!")
			},
			err: "failed on the 'max' tag",
		},
		{
			name: "too short",
			modifier: func(up *model.UpdateProject) {
				up.Name = aws.String("t")
			},
			err: "failed on the 'min' tag",
		},
		{
			name: "too long",
			modifier: func(up *model.UpdateProject) {
				up.Description = aws.String("This text represents a project description that's unfortunately not valid") // over by 1 char
			},
			err: "failed on the 'max' tag",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			up := model.UpdateProject{}

			tc.modifier(&up)

			err := up.Validate()

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

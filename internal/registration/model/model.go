// Package model contains domain models for the application.
package model

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// NewValidator returns a new validator aware of json tags.
func NewValidator() *validator.Validate {
	v := validator.New()
	v.RegisterTagNameFunc(jsonTagName)
	return v
}

func jsonTagName(fld reflect.StructField) string {
	name := strings.SplitN(fld.Tag.Get("json"), ",", 1)[0]
	if name == "-" {
		return ""
	}
	return name
}

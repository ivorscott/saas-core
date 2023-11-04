// Package main runs the service.
package main

import (
	"github.com/devpies/saas-core/internal/registration"
)

func main() {
	err := registration.Run()
	if err != nil {
		panic(err)
	}
}

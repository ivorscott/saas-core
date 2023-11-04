// Package main runs the service.
package main

import (
	"github.com/devpies/saas-core/internal/tenant"
)

func main() {
	err := tenant.Run()
	if err != nil {
		panic(err)
	}
}

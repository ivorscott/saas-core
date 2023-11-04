// Package main runs the service.
package main

import (
	"github.com/devpies/saas-core/internal/billing"
)

func main() {
	err := billing.Run()
	if err != nil {
		panic(err)
	}
}

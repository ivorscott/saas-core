// Package main runs the service.
package main

import "github.com/devpies/saas-core/internal/subscription"

func main() {
	err := subscription.Run()
	if err != nil {
		panic(err)
	}
}

package main

import "github.com/devpies/saas-core/internal/identity"

func main() {
	err := identity.Run()
	if err != nil {
		panic(err)
	}
}

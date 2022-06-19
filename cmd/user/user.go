package main

import "github.com/devpies/saas-core/internal/user"

func main() {
	err := user.Run()
	if err != nil {
		panic(err)
	}
}

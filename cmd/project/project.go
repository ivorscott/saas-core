package main

import (
	"github.com/devpies/saas-core/internal/project"
)

func main() {
	err := project.Run()
	if err != nil {
		panic(err)
	}
}

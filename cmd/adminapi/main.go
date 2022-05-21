package main

import "github.com/devpies/core/internal/adminapi"

func main() {
	err := adminapi.Run()
	if err != nil {
		panic(err)
	}
}

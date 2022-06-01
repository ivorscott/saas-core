package main

import (
	"embed"

	"github.com/devpies/core/internal/adminclient"
)

//go:embed static
var staticFS embed.FS

func main() {
	err := adminclient.Run(staticFS)
	if err != nil {
		panic(err)
	}
}

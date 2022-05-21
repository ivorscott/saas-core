package main

import (
	"embed"
	app "github.com/devpies/core/internal/admin-client"
)

//go:embed static
var staticFS embed.FS

func main() {
	err := app.Run(staticFS)
	if err != nil {
		panic(err)
	}
}

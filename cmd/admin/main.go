// Package main runs the service.
package main

import (
	"embed"

	"github.com/devpies/saas-core/internal/admin"
)

//go:embed static
var staticFS embed.FS

func main() {
	err := admin.Run(staticFS)
	if err != nil {
		panic(err)
	}
}

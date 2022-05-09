package handlers

import (
	"net/http"
	"strings"

	"github.com/rs/cors"
)

// Cors sets up a new Cors Handler
func Cors(origins string) *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   ParseCorsOrigins(origins),
		AllowedHeaders:   []string{"Authorization", "Cache-Control", "Content-Type", "Strict-Transport-Security"},
		AllowedMethods:   []string{http.MethodOptions, http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut, http.MethodPatch},
		AllowCredentials: true,
	})
}

// ParseCorsOrigins splits a comma separated origins string into multiple allowed origins
func ParseCorsOrigins(origins string) []string {
	rawOrigins := strings.Split(origins, ",")
	o := make([]string, 0, len(rawOrigins))

	for _, v := range rawOrigins {
		trimmed := strings.TrimSpace(v)
		o = append(o, trimmed)
	}
	return o
}

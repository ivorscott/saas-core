package handlers

import (
	"github.com/rs/cors"
	"net/http"
	"strings"
)

func Cors(origins string) *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   ParseOrigins(origins),
		AllowedHeaders:   []string{"Authorization", "Cache-Control", "Content-Type", "Strict-Transport-Security"},
		AllowedMethods:   []string{http.MethodOptions, http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut, http.MethodPatch},
		AllowCredentials: true,
	})
}

func ParseOrigins(origins string) []string {
	rawOrigins := strings.Split(origins, ",")
	o := make([]string, 0, len(rawOrigins))

	for _, v := range rawOrigins {
		trimmed := strings.TrimSpace(v)
		o = append(o, trimmed)
	}
	return o
}
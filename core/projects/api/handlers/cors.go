package handlers

import (
	"github.com/rs/cors"
	"net/http"
	"strings"
)

func Cors(origins string) *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   parseOrigins(origins),
		AllowedHeaders:   []string{"Authorization", "Cache-Control", "Content-Type", "Strict-Transport-Security"},
		AllowedMethods:   []string{http.MethodOptions, http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut, http.MethodPatch},
		AllowCredentials: true,
	})
}

func parseOrigins(origins string) []string {
	rawOrigins := strings.Split(origins, ",")
	o := make([]string, len(rawOrigins))

	for _, v := range rawOrigins {
		o = append(o, strings.TrimSpace(v))
	}
	return o
}
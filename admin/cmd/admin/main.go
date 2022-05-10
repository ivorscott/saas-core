package main

import (
	"fmt"
	tm "github.com/devpies/core/tenant/pkg/model"
	"html/template"
	"net/http"

	"github.com/devpies/core/pkg/log"
)

func index(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"templates/layout.gohtml",
		"templates/index.gohtml",
	}

	m := tm.Tenant{
		ID:   "2134",
		Name: "test",
		URL:  "url-here",
	}

	templates := template.Must(template.ParseFiles(files...))
	templates.ExecuteTemplate(w, "layout", struct{}{})
	fmt.Println(m, r.URL)
}

var logPath = "../.log/out.log"

func main() {
	var err error

	logger, Sync := log.NewLoggerOrPanic(logPath)
	defer Sync()

	mux := http.NewServeMux()
	files := http.FileServer(http.Dir("public"))
	mux.Handle("/static/", http.StripPrefix("/static/", files))
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("/", index)

	server := &http.Server{
		Addr:    "0.0.0.0:8989",
		Handler: mux,
	}

	if err = server.ListenAndServe(); err != nil {
		logger.Fatal(err.Error())
	}
}

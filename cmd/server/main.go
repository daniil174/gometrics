package main

import (
	"github.com/daniil174/gometrics/cmd/server/handlers"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	//mux := http.NewServeMux()
	// mux.HandleFunc(`/`, handlers.MainPage)
	//r.Routes("/", func(r chi.Router) {
	//	r.Get("/", handlers.MainPage)
	//})

	r.Get("/", handlers.MainPage)
	r.Post("/update/{type}/{name}/{value}", handlers.UpdateMetrics)
	r.Get("/value/{type}/{name}", handlers.GetMetric)

	//mux.HandleFunc(`/update/{type}/{name}/{value}`, handlers.UpdateMetrics)

	err := http.ListenAndServe(`:8080`, r)
	if err != nil {
		panic(err)
	}
}

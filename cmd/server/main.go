package main

import (
	"net/http"

	"github.com/daniil174/gometrics/cmd/server/handlers"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	r.Get("/", handlers.MainPage)
	r.Post("/update/{type}/{name}/{value}", handlers.UpdateMetrics)
	r.Get("/value/{type}/{name}", handlers.GetMetric)

	err := http.ListenAndServe(`:8080`, r)
	if err != nil {
		panic(err)
	}
}

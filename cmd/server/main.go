package main

import (
	"flag"
	"net/http"

	"github.com/daniil174/gometrics/cmd/server/handlers"
	"github.com/go-chi/chi/v5"
)

func main() {
	serverAddr := flag.String("a", "localhost:8080", "server address and port, example 127.0.0.1:8080")
	flag.Parse()

	r := chi.NewRouter()
	r.Get("/", handlers.MainPage)
	r.Post("/update/{type}/{name}/{value}", handlers.UpdateMetrics)
	r.Get("/value/{type}/{name}", handlers.GetMetric)

	err := http.ListenAndServe(*serverAddr, r)
	if err != nil {
		panic(err)
	}
}

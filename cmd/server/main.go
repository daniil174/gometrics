package main

import (
	"github.com/daniil174/gometrics/cmd/server/handlers"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	//mux.HandleFunc(`/`, handlers.MainPage)
	mux.HandleFunc(`/update/{type}/{name}/{value}`, handlers.UpdateMetrics)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}

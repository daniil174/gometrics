package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/caarlos0/env/v11"
	"github.com/daniil174/gometrics/cmd/server/handlers"
	"github.com/go-chi/chi/v5"
)

type Config struct {
	Addr string `env:"ADDRESS"`
}

var serverAddr string

func ConfigFromEnv() error {
	cfg, errConf := env.ParseAs[Config]()
	if errConf != nil {
		return errConf
	}
	fmt.Printf("ADDRESS=%s=", cfg.Addr)
	serverAddr = cfg.Addr
	if serverAddr == "" {
		flag.StringVar(&serverAddr, "a", "localhost:8080", "server address and port, example 127.0.0.1:8080")
		flag.Parse()
	}
	fmt.Printf("serverAddr=%s=", serverAddr)
	return nil
}

func main() {
	_ = ConfigFromEnv()
	r := chi.NewRouter()
	r.Get("/", handlers.MainPage)
	r.Post("/update/{type}/{name}/{value}", handlers.UpdateMetrics)
	r.Get("/value/{type}/{name}", handlers.GetMetric)

	err := http.ListenAndServe(serverAddr, r)
	if err != nil {
		panic(err)
	}
}

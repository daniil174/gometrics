package main

import (
	"log"
	"net/http"
	"time"

	"github.com/daniil174/gometrics/cmd/server/compress"
	"github.com/daniil174/gometrics/cmd/server/handlers"
	"github.com/daniil174/gometrics/cmd/server/servconfig"
	"github.com/daniil174/gometrics/cmd/server/servlogger"
	"github.com/go-chi/chi/v5"
)

func main() {

	tmpCfg, _ := servconfig.SetConfig()

	logger := servlogger.Sugar

	r := chi.NewRouter()
	r.Use(servlogger.AddLogging)
	r.Use(compress.GzipHandleEncode)
	r.Use(compress.GzipHandleDecode)

	//mem := storage.NewMemStorage()
	//stor := storage.New()
	//stor.OpenFile(*tmpCfg)

	//defer stor.CloseFile()

	if tmpCfg.Restore {
		err := handlers.MemStrg.ReadFile(tmpCfg.FileStoragePath)
		if err != nil {
			servlogger.Sugar.Fatalw(err.Error(), "event", "on load metrics from file")
		}
	}

	go func() {
		for {
			interval := time.Duration(tmpCfg.StoreInterval) * time.Second
			// if interval == 0 {
			// 	interval = 100 * time.Microsecond // Установите разумное значение по умолчанию
			// }
			time.Sleep(interval)
			//stor.SaveMetricsToFile()
			err := handlers.MemStrg.SaveMetricsToFile(tmpCfg.FileStoragePath)
			if err != nil {
				servlogger.Sugar.Fatalw(err.Error(), "event", "on save metrics in file")
			}
			log.Printf("event", "save metrics in file", interval)
		}
	}()

	//defer handlers.MemStrg.CloseFile()

	// Добавляем просмотр логов по запросу "http://serverAddr/logs"
	// Временно убрал, из-за автотестов
	// r.Get("/*", servlogger.Logs)

	r.Get("/", handlers.MainPage)

	r.Post("/update/{type}/{name}/{value}", handlers.UpdateMetrics)
	r.Get("/value/{type}/{name}", handlers.GetMetric)

	// Роуты для JSON запросов
	r.Post("/update/", handlers.UpdateMetrics2)
	r.Post("/value/", handlers.GetMetric2)

	err2 := http.ListenAndServe(tmpCfg.ServerAddress, r)
	if err2 != nil {
		panic(err2)
		logger.Fatalw(err2.Error(), "event", "on start server")
	}

}

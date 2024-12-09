package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/daniil174/gometrics/internal/server/compress"
	"github.com/daniil174/gometrics/internal/server/handlers"
	"github.com/daniil174/gometrics/internal/server/servconfig"
	"github.com/daniil174/gometrics/internal/server/servlogger"
	"github.com/daniil174/gometrics/internal/storage"

	"github.com/go-chi/chi/v5"
)

func main() {

	tmpCfg, _ := servconfig.SetConfig()
	logger := servlogger.Sugar

	// Если есть база данных - берем всё из неё, если нет, то из файла, если и его нет, то всё в памяти храним
	if tmpCfg.DatabaseDsn != "" {
		storage.MemStrg.SetMemType(storage.DB)
		storage.StartDB(tmpCfg.DatabaseDsn)

		defer storage.CloseDB()

		if !tmpCfg.Restore {
			storage.MemStrg.ResetDBandSetZeroValue()
		}
		storage.MemStrg.ReadMetricFromDB()

		interval := time.Duration(tmpCfg.StoreInterval) * time.Second
		go func() {
			for {
				// if interval == 0 {
				// 	interval = 100 * time.Microsecond // Установите разумное значение по умолчанию
				// }
				time.Sleep(interval)
				err := storage.MemStrg.WriteMetricToDB()
				if err != nil {
					servlogger.Sugar.Infow(err.Error(), "event", "on save metrics in file")
				}
				//log.Printf("event", "save metrics in file", interval)
			}
		}()

	} else if tmpCfg.FileStoragePath != "" {
		storage.MemStrg.SetMemType(storage.FILE)

		// если флаг Restore=true - пробуем восстановить из файла
		if tmpCfg.Restore {
			fmt.Printf("\n Try to read data from file %s\n", tmpCfg.FileStoragePath)
			err := storage.MemStrg.ReadFile(tmpCfg.FileStoragePath)
			if err != nil {
				servlogger.Sugar.Fatalw(err.Error(), "event", "on load metrics from file")
			}

		}

		// используя интервал в секундах - пишем данные в файл
		interval := time.Duration(tmpCfg.StoreInterval) * time.Second
		go func() {
			for {
				// if interval == 0 {
				// 	interval = 100 * time.Microsecond // Установите разумное значение по умолчанию
				// }
				time.Sleep(interval)
				err := storage.MemStrg.SaveMetricsToFile(tmpCfg.FileStoragePath)
				if err != nil {
					servlogger.Sugar.Infow(err.Error(), "event", "on save metrics in file")
				}
				//log.Printf("event", "save metrics in file", interval)
			}
		}()
	} else {
		storage.MemStrg.SetMemType(storage.NONE)
	}

	r := chi.NewRouter()
	r.Use(servlogger.AddLogging)
	r.Use(compress.GzipHandleEncode)
	r.Use(compress.GzipHandleDecode)

	//mem := storage.NewMemStorage()
	//stor := storage.New()
	//stor.OpenFile(*tmpCfg)

	//defer stor.CloseFile()

	// Добавляем просмотр логов по запросу "http://serverAddr/logs"
	// Временно убрал, из-за автотестов
	// r.Get("/*", servlogger.Logs)

	r.Get("/", handlers.MainPage)
	r.Get("/ping", handlers.DBhealthcheck)

	r.Post("/update/{type}/{name}/{value}", handlers.UpdateMetrics)
	r.Get("/value/{type}/{name}", handlers.GetMetric)

	// Роуты для JSON запросов
	r.Post("/update/", handlers.UpdateMetrics2)
	r.Post("/value/", handlers.GetMetric2)

	err := http.ListenAndServe(tmpCfg.ServerAddress, r)
	if err != nil {
		//panic(err2)
		logger.Fatalw(err.Error(), "event", "on start server")
	}

}

package main

import (
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/daniil174/gometrics/cmd/server/handlers"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var sugar zap.SugaredLogger

const (
	LogMaxSize    = 10 // megabytes
	LogMaxBackups = 3
	LogMaxAge     = 7 // days
)

func createLogger() *zap.Logger {
	stdout := zapcore.AddSync(os.Stdout)

	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/logs",
		MaxSize:    LogMaxSize, // megabytes
		MaxBackups: LogMaxBackups,
		MaxAge:     LogMaxAge, // days
	})

	level := zap.NewAtomicLevelAt(zap.InfoLevel)

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, level),
		zapcore.NewCore(fileEncoder, file, level),
	)

	return zap.New(core)
}

type Config struct {
	Addr string `env:"ADDRESS"`
}

var serverAddr string

func ConfigFromEnv() error {
	cfg, errConf := env.ParseAs[Config]()
	if errConf != nil {
		return errConf
	}
	fmt.Printf("ADDRESS=%s=\n", cfg.Addr)
	serverAddr = cfg.Addr
	if serverAddr == "" {
		flag.StringVar(&serverAddr, "a", "localhost:8080", "server address and port, example 127.0.0.1:8080")
		flag.Parse()
	}
	fmt.Printf("serverAddr=%s=\n", serverAddr)
	return nil
}

type (
	// Берём структуру для хранения сведений об ответе
	responseData struct {
		status int
		size   int
	}

	// Добавляем реализацию http.ResponseWriter
	loggingResponseWriter struct {
		http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
		responseData        *responseData
	}
)

// Write для логирования размера ответа
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	// Записываем ответ, используя оригинальный http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // захватываем размер
	return size, err
}

// WriteHeader Для логирования кода статуса
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// Записываем код статуса, используя оригинальный http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // захватываем код статуса
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	return w.Writer.Write(b)
}

func gzipHandleDecode(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем заголовок Content-Encoding
		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			// если gzip не использовался, передаём управление
			// дальше без изменений
			next.ServeHTTP(w, r)
			return
		}

		// создаём *gzip.Reader, который будет читать тело запроса
		// и распаковывать его
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer gz.Close()

		body, err := io.ReadAll(gz)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(body))
		next.ServeHTTP(w, r)
	})
}

func gzipHandleEncode(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// проверяем, что клиент поддерживает gzip-сжатие
		// это упрощённый пример. В реальном приложении следует проверять все
		// значения r.Header.Values("Accept-Encoding") и разбирать строку
		// на составные части, чтобы избежать неожиданных результатов
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			// если gzip не поддерживается, передаём управление
			// дальше без изменений
			next.ServeHTTP(w, r)
			return
		}

		// создаём gzip.Writer поверх текущего w
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		// передаём обработчику страницы переменную типа gzipWriter для вывода данных
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

func WithLogging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := createLogger()
		defer logger.Sync()

		sugar := *logger.Sugar()

		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}

		lw := loggingResponseWriter{
			ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
			responseData:   responseData,
		}
		h.ServeHTTP(&lw, r) // внедряем реализацию http.ResponseWriter

		duration := time.Since(start)

		sugar.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"duration", duration,
			"status", responseData.status, // получаем перехваченный код статуса ответа
			"size", responseData.size, // получаем перехваченный размер ответа
		)
	})
}

func main() {
	_ = ConfigFromEnv()
	r := chi.NewRouter()
	r.Use(WithLogging)
	r.Use(gzipHandleEncode)
	r.Use(gzipHandleDecode)

	// Добавляем просмотр логов по запросу "http://serverAddr/log"
	// r.Get("/*", handlers.Logs)

	r.Get("/", handlers.MainPage)

	r.Post("/update/{type}/{name}/{value}", handlers.UpdateMetrics)
	r.Post("/update/", handlers.UpdateMetrics2)
	r.Post("/value/", handlers.GetMetric2)
	r.Get("/value/{type}/{name}", handlers.GetMetric)

	err2 := http.ListenAndServe(serverAddr, r)
	if err2 != nil {
		sugar.Fatalw(err2.Error(), "event", "on start server")
	}
}

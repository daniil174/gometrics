package handlers

import (
	"errors"
	"fmt"
	"github.com/daniil174/gometrics/cmd/server/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
	"sort"
	"strconv"
)

var m = storage.NewMemStorage()

func UpdateMetrics(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")

	body := fmt.Sprintf("metricType : %s\n metricName : %s\n ", metricType, metricName)

	if metricName == "" {
		http.Error(w, "Metric must have NAME", http.StatusNotFound)
		return
	}

	switch metricType {
	case "counter":
		{
			metricValue, err := strconv.ParseInt(chi.URLParam(r, "value"), 10, 64)
			if err != nil {
				http.Error(w, "Metric counter must have INT value", http.StatusBadRequest)
				return
			}

			err = m.AddCounter(metricName, metricValue)
			if err != nil {
				if errors.Is(err, storage.ErrMetricDidntExist) {
					http.Error(w, "Metric counter did't exists", http.StatusNotFound)
				} else {
					return
				}
			}

			res, _ := m.GetCounter(metricName)
			body += fmt.Sprintf("metricValue : %d\n", res)
		}
	case "gauge":
		{
			metricValue, err := strconv.ParseFloat(chi.URLParam(r, "value"), 64)
			if err != nil {
				http.Error(w, "Metric gauge must have float64 value", http.StatusBadRequest)
				return
			}

			err = m.RewriteGauge(metricName, metricValue)
			if err != nil {
				if errors.Is(err, storage.ErrMetricDidntExist) {
					http.Error(w, "Metric counter did't exists", http.StatusNotFound)
				} else {
					return
				}
			}
			res, _ := m.GetGauge(metricName)
			body += fmt.Sprintf("metricValue : %f\n", res)
		}
	default:
		{
			http.Error(w, "Metric TYPE must be 'counter' or 'gauge'", http.StatusBadRequest)
			return
		}
	}

	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(body))
	if err != nil {
		return
	}
}

func GetMetric(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")

	switch metricType {
	case "counter":
		{
			resp, err := m.GetCounter(metricName)
			if err != nil {
				if errors.Is(err, storage.ErrMetricDidntExist) {
					http.Error(w, "Metric did't exists", http.StatusNotFound)
				} else {
					return
				}
			}

			_, err = fmt.Fprintf(w, "%d", resp)
			if err != nil {
				return
			}
		}
	case "gauge":
		{
			resp, err := m.GetGauge(metricName)
			if err != nil {
				if errors.Is(err, storage.ErrMetricDidntExist) {
					http.Error(w, "Metric did't exists", http.StatusNotFound)
				} else {
					return
				}
			}

			_, err = fmt.Fprintf(w, "%s", strconv.FormatFloat(resp, 'f', -1, 64))
			if err != nil {
				return
			}
		}
	default:
		{
			http.Error(w, "Metric TYPE must be 'counter' or 'gauge'", http.StatusBadRequest)
			return
		}
	}

	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusOK)
}

func MainPage(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusOK)
	var body = ""
	for n, v := range m.Counter {
		body += fmt.Sprintf("Metric name: %s = %d \n", n, v)
	}

	// Sort Gauge metrics by name
	keys := make([]string, 0, len(m.Gauge))
	for k := range m.Gauge {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		body += fmt.Sprintf("Metric name: %s = %s \n", k, strconv.FormatFloat(m.Gauge[k], 'f', -1, 64))
	}

	_, err := w.Write([]byte(body))
	if err != nil {
		return
	}
}

/* func Logs(w http.ResponseWriter, r *http.Request) {
	// даем загружать только файлы из папки "logs"
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "logs"))

	rctx := chi.RouteContext(r.Context())
	pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
	fs := http.StripPrefix(pathPrefix, http.FileServer(filesDir))
	fs.ServeHTTP(w, r)
} */

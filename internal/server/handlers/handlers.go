package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/daniil174/gometrics/internal/storage"
	"github.com/go-chi/chi/v5"
)

var MemStrg = storage.NewMemStorage()

func UpdateMetrics2(w http.ResponseWriter, r *http.Request) {
	var metric storage.Metrics
	w.Header().Set("content-type", "application/json")

	jsDec := json.NewDecoder(r.Body)
	if err := jsDec.Decode(&metric); err != nil {
		http.Error(w, "Post request must have body", http.StatusInternalServerError)
		return
	}

	switch metric.MType {
	case "gauge":
		{
			err := MemStrg.RewriteGauge(metric.ID, *metric.Value)
			if err != nil {
				if errors.Is(err, storage.ErrMetricDidntExist) {
					http.Error(w, "Metric counter did't exists", http.StatusNotFound)
				} else {
					return
				}
			}
			*metric.Value, _ = MemStrg.GetGauge(metric.ID)
			err = json.NewEncoder(w).Encode(metric)
			if err != nil {
				return
			}

		}
	case "counter":
		{
			err := MemStrg.AddCounter(metric.ID, *metric.Delta)
			if err != nil {
				if errors.Is(err, storage.ErrMetricDidntExist) {
					http.Error(w, "Metric counter did't exists", http.StatusNotFound)
				} else {
					return
				}
			}

			*metric.Delta, _ = MemStrg.GetCounter(metric.ID)
			err = json.NewEncoder(w).Encode(metric)
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

	w.WriteHeader(http.StatusOK)
}

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

			err = MemStrg.AddCounter(metricName, metricValue)
			if err != nil {
				if errors.Is(err, storage.ErrMetricDidntExist) {
					http.Error(w, "Metric counter did't exists", http.StatusNotFound)
				} else {
					return
				}
			}

			res, _ := MemStrg.GetCounter(metricName)
			body += fmt.Sprintf("metricValue : %d\n", res)
		}
	case "gauge":
		{
			metricValue, err := strconv.ParseFloat(chi.URLParam(r, "value"), 64)
			if err != nil {
				http.Error(w, "Metric gauge must have float64 value", http.StatusBadRequest)
				return
			}

			err = MemStrg.RewriteGauge(metricName, metricValue)
			if err != nil {
				if errors.Is(err, storage.ErrMetricDidntExist) {
					http.Error(w, "Metric counter did't exists", http.StatusNotFound)
				} else {
					return
				}
			}
			res, _ := MemStrg.GetGauge(metricName)
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

func GetMetric2(w http.ResponseWriter, r *http.Request) {

	//if r.Header.Get("Content-Type") != "application/json" {
	//	http.Error(w, "Wrong content type ", http.StatusInternalServerError)
	//	return
	//}

	var metric storage.Metrics

	jsDec := json.NewDecoder(r.Body)
	if err := jsDec.Decode(&metric); err != nil {
		http.Error(w, "Post request must have body", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")

	switch metric.MType {
	case "gauge":
		{
			resp, err := MemStrg.GetGauge(metric.ID)
			if err != nil {
				if errors.Is(err, storage.ErrMetricDidntExist) {
					http.Error(w, "Metric did't exists", http.StatusNotFound)
					return
				} else {
					return
				}
			}

			respMetric := storage.Metrics{
				ID:    metric.ID,
				MType: metric.MType,
				Value: &resp,
			}
			jsEnc := json.NewEncoder(w)
			if err := jsEnc.Encode(respMetric); err != nil {
				http.Error(w, "Metric value problem", http.StatusInternalServerError)
			}
		}
	case "counter":
		{
			resp, err := MemStrg.GetCounter(metric.ID)
			if err != nil {
				if errors.Is(err, storage.ErrMetricDidntExist) {
					http.Error(w, "Metric did't exists", http.StatusNotFound)
				}
				return

			}

			respMetric := storage.Metrics{
				ID:    metric.ID,
				MType: metric.MType,
				Delta: &resp,
			}
			jsEnc := json.NewEncoder(w)
			if err := jsEnc.Encode(respMetric); err != nil {
				http.Error(w, "Metric value problem", http.StatusInternalServerError)
				return
			}
		}
	default:
		{
			http.Error(w, "Metric TYPE must be 'counter' or 'gauge'", http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusOK)

}

func GetMetric(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")

	switch metricType {
	case "counter":
		{
			resp, err := MemStrg.GetCounter(metricName)
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
			resp, err := MemStrg.GetGauge(metricName)
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
	w.Header().Set("content-type", "text/html")
	w.WriteHeader(http.StatusOK)
	var body = ""
	for n, v := range MemStrg.Counter {
		body += fmt.Sprintf("<br> Metric name: %s = %d \n", n, v)
	}

	// Sort Gauge metrics by name
	keys := make([]string, 0, len(MemStrg.Gauge))
	for k := range MemStrg.Gauge {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		body += fmt.Sprintf("<br> Metric name: %s = %s \n ", k, strconv.FormatFloat(MemStrg.Gauge[k], 'f', -1, 64))
	}

	_, err := w.Write([]byte(body))
	if err != nil {
		return
	}
}

func DBhealthcheck(w http.ResponseWriter, _ *http.Request) {
	//w.Header().Set("content-type", "text/html")
	err := storage.PgDataBase.Ping()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	//if _, err := storage.PingDB(); err != nil {
	//	w.WriteHeader(http.StatusOK)
	//	return
	//}
	//w.WriteHeader(http.StatusInternalServerError)

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

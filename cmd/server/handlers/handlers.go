package handlers

import (
	"fmt"
	"github.com/daniil174/gometrics/cmd/server/storage"
	"net/http"
	"strconv"
	"strings"
)

var m = storage.NewMemStorage()

func UpdateMetrics(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	metricType := r.PathValue("type")
	metricName := r.PathValue("name")
	body := fmt.Sprintf("metricType : %s\n metricName : %s\n ", metricType, metricName)

	if metricName == "" {
		http.Error(w, "Metric must have NAME", http.StatusNotFound)
		return
	}

	switch metricType {
	case "counter":
		{

			metricValue, err := strconv.ParseInt(r.PathValue("value"), 10, 64)
			if err != nil {
				http.Error(w, "Metric counter must have INT value", http.StatusBadRequest)
				return
			}
			if _, ok := m.Counter[metricName]; ok {
				err := m.AddCounter(metricName, metricValue)
				if err != nil {
					return
				}
				//m.counter[metricName] += metricValue
			} else {
				m.Counter[metricName] = metricValue
			}
			body += fmt.Sprintf("metricValue : %d\n", m.Counter[metricName])
		}
	case "gauge":
		{
			metricValue, err := strconv.ParseFloat(strings.TrimSpace(r.PathValue("value")), 64)
			if err != nil {
				// ... handle error
				http.Error(w, "Metric gauge must have float64 value", http.StatusBadRequest)
				return
			}
			if _, ok := m.Gauge[metricName]; ok {
				err := m.RewriteGauge(metricName, metricValue)
				if err != nil {
					return
				}
				//m.gauge[metricName] = metricValue
			} else {
				m.Gauge[metricName] = metricValue
			}
			body += fmt.Sprintf("metricValue : %f\n", m.Gauge[metricName])
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

func MainPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./README.md")
}

package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type ()

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
}

func (m *MemStorage) addCounter(name string, value int64) error {
	m.counter[name] += value
	return nil
}

func (m *MemStorage) rewriteGauge(name string, value float64) error {
	m.gauge[name] = value
	return nil
}

var m = NewMemStorage()

func updateMetrics(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	} else {
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
					// ... handle error
					http.Error(w, "Metric counter must have INT value", http.StatusBadRequest)
					return
				}
				if _, ok := m.counter[metricName]; ok {
					m.addCounter(metricName, metricValue)
					//m.counter[metricName] += metricValue
				} else {
					m.counter[metricName] = metricValue
				}
				body += fmt.Sprintf("metricValue : %d\n", m.counter[metricName])
			}
		case "gauge":
			{
				metricValue, err := strconv.ParseFloat(strings.TrimSpace(r.PathValue("value")), 64)
				if err != nil {
					// ... handle error
					http.Error(w, "Metric gauge must have float64 value", http.StatusBadRequest)
					return
				}
				if _, ok := m.gauge[metricName]; ok {
					m.rewriteGauge(metricName, metricValue)
					//m.gauge[metricName] = metricValue
				} else {
					m.gauge[metricName] = metricValue
				}
				body += fmt.Sprintf("metricValue : %f\n", m.gauge[metricName])
			}
		default:
			{
				http.Error(w, "Metric TYPE must be 'counter' or 'gauge'", http.StatusBadRequest)
				return
			}
		}

		w.Header().Set("content-type", "text/plain")
		//w.Header().Set("content-type", "charset=utf-8")
		w.WriteHeader(http.StatusOK)
		//_, err := w.Write([]byte(body))
		//if err != nil {
		//	return
		//}
	}

	//if r.Method == http.MethodPost {
	//	login := r.FormValue("login")
	//	password := r.FormValue("password")
	//	if Auth(login, password) {
	//		io.WriteString(w, "Добро пожаловать!")
	//	} else {
	//		http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
	//	}
	//	return
	//} else {
	//	io.WriteString(w, form)
	//}
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./README.md")
}

func main() {

	mux := http.NewServeMux()
	//mux.HandleFunc(`/`, mainPage)
	mux.HandleFunc(`/update/{type}/{name}/{value}`, updateMetrics)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}

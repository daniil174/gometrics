package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"
)

type GaugeMetric struct {
	Name string
	//Type  string
	Value float64
}

var m runtime.MemStats

func CollectGaugeMetrics() []GaugeMetric {

	runtime.ReadMemStats(&m)

	return []GaugeMetric{
		{Name: "Alloc", Value: float64(m.Alloc)},
		{Name: "BuckHashSys", Value: float64(m.BuckHashSys)},
		{Name: "Frees", Value: float64(m.Frees)},
		{Name: "GCCPUFraction", Value: m.GCCPUFraction},
		{Name: "GCSys", Value: float64(m.GCSys)},
		{Name: "HeapAlloc", Value: float64(m.HeapAlloc)},
		{Name: "HeapIdle", Value: float64(m.HeapIdle)},
		{Name: "HeapInuse", Value: float64(m.HeapInuse)},
		{Name: "HeapObjects", Value: float64(m.HeapObjects)},
		{Name: "HeapReleased", Value: float64(m.HeapReleased)},
		{Name: "HeapSys", Value: float64(m.HeapSys)},
		{Name: "LastGC", Value: float64(m.LastGC)},
		{Name: "Lookups", Value: float64(m.Lookups)},
		{Name: "MCacheInuse", Value: float64(m.MCacheInuse)},
		{Name: "MCacheSys", Value: float64(m.MCacheSys)},
		{Name: "MSpanInuse", Value: float64(m.MSpanInuse)},
		{Name: "MSpanSys", Value: float64(m.MSpanSys)},
		{Name: "Mallocs", Value: float64(m.Mallocs)},
		{Name: "NextGC", Value: float64(m.NextGC)},
		{Name: "NumForcedGC", Value: float64(m.NumForcedGC)},
		{Name: "NumGC", Value: float64(m.NumGC)},
		{Name: "OtherSys", Value: float64(m.OtherSys)},
		{Name: "PauseTotalNs", Value: float64(m.PauseTotalNs)},
		{Name: "StackInuse", Value: float64(m.StackInuse)},
		{Name: "StackSys", Value: float64(m.StackSys)},
		{Name: "Sys", Value: float64(m.Sys)},
		{Name: "TotalAlloc", Value: float64(m.TotalAlloc)},
		//{Name: "PollCount", Type: "counter", Value: int64(pollCount)},
		{Name: "RandomValue", Value: rand.Float64() * 100},
	}
}

func SendMetrics(pollCount int64) {
	for _, v := range CollectGaugeMetrics() {
		URL := "http://localhost:8080/update/gauge/" + v.Name + "/" + strconv.FormatFloat(v.Value, 'f', -1, 64)
		res, err := http.Post(URL, "Content-Type: text/plain", nil)
		if err != nil {
			fmt.Printf("error making http request: %s\n", err)
			os.Exit(1)
		}
		//fmt.Printf("client: got response!\n")
		fmt.Printf("client: status code: %d\n", res.StatusCode)

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				fmt.Printf("error: %s\n", err)
				os.Exit(1)
			}
		}(res.Body)
		//if res.StatusCode == http.StatusOK {
		//	bodyBytes, err := io.ReadAll(res.Body)
		//	if err != nil {
		//		log.Fatal(err)
		//	}
		//	bodyString := string(bodyBytes)
		//	fmt.Printf("client: status body:", bodyString)
		//}
	}

	URL := "http://localhost:8080/update/gauge/PollCount/" + strconv.FormatInt(pollCount, 10)
	res, err := http.Post(URL, "Content-Type: text/plain", nil)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("client: status code: %d\n", res.StatusCode)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("error: %s\n", err)
			os.Exit(1)
		}
	}(res.Body)
	//if res.StatusCode == http.StatusOK {
	//	bodyBytes, err := io.ReadAll(res.Body)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	bodyString := string(bodyBytes)
	//	fmt.Printf("client: status body:", bodyString)
	//}
}

var pollInterval = 2 * time.Second
var reportInterval = 10 * time.Second

var pullCount int64 = 0

func main() {

	//countTime := 0

	startTimePoll := time.Now()
	startTimeReport := time.Now()

	for {
		time.Sleep(time.Second)
		finishTimePoll := time.Now()
		finishTimeReport := time.Now()
		//countTime++
		//fmt.Printf("Time: %d\n", countTime)

		if finishTimePoll.Sub(startTimePoll) >= pollInterval {
			CollectGaugeMetrics()
			pullCount++
			//fmt.Printf("Pull Count: %d\n", countTime)
			startTimePoll = time.Now()
		}

		if finishTimeReport.Sub(startTimeReport) >= reportInterval {
			SendMetrics(pullCount)
			//fmt.Printf("Send Report: %d\n", countTime)
			startTimeReport = time.Now()

		}

	}
}

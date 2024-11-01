package main

import (
	"fmt"
	"github.com/daniil174/gometrics/internal/memstats"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

func SendMetrics2() {
	client := resty.New()
	for _, v := range memstats.CollectGaugeMetrics() {
		_, err := client.R().SetPathParams(map[string]string{
			"Name":  v.Name,
			"Value": fmt.Sprintf("%f", v.Value),
		}).Post("http://localhost:8080/update/gauge/{Name}/{Value}")
		if err != nil {
			panic(err)
		}
		println("gauge ok")
	}

	_, err := client.R().SetPathParams(map[string]string{
		"Name":  "PollCount",
		"Value": strconv.FormatInt(memstats.PullCount+1, 10),
	}).Post("http://localhost:8080/update/counter/{Name}/{Value}")

	if err != nil {
		panic(err)
	}

	fmt.Printf("counter ok %d", memstats.PullCount+1)
}

var pollInterval = 2 * time.Second
var reportInterval = 10 * time.Second

func main() {
	// countTime := 0
	startTimePoll := time.Now()
	startTimeReport := time.Now()

	for {
		time.Sleep(time.Second)
		finishTimePoll := time.Now()
		finishTimeReport := time.Now()
		// countTime++
		// fmt.Printf("Time: %d\n", countTime)
		if finishTimePoll.Sub(startTimePoll) >= pollInterval {
			memstats.CollectGaugeMetrics()
			//pullCount++
			// fmt.Printf("Pull Count: %d\n", countTime)
			startTimePoll = time.Now()
		}
		if finishTimeReport.Sub(startTimeReport) >= reportInterval {
			SendMetrics2()
			// fmt.Printf("Send Report: %d\n", countTime)
			startTimeReport = time.Now()
		}
	}
}

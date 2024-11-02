package main

import (
	"flag"
	"fmt"
	"strconv"
	"time"

	"github.com/daniil174/gometrics/internal/memstats"

	"github.com/go-resty/resty/v2"
)

const NanoSecToSec = 1000 * 1000 * 1000
const defaultPollInterval = 2
const defaultReportInterval = 10

var pollInterval int
var reportInterval int
var serverAddr string

func SendMetrics2() {
	client := resty.New()
	for _, v := range memstats.CollectGaugeMetrics() {
		_, err := client.R().SetPathParams(map[string]string{
			"serverAddressAndPort": serverAddr,
			"Name":                 v.Name,
			"Value":                fmt.Sprintf("%f", v.Value),
		}).Post("http://{serverAddressAndPort}/update/gauge/{Name}/{Value}")
		if err != nil {
			panic(err)
		}
		println("gauge ok")
	}

	_, err := client.R().SetPathParams(map[string]string{
		"serverAddressAndPort": serverAddr,
		"Name":                 "PollCount",
		"Value":                strconv.FormatInt(memstats.PullCount+1, 10),
	}).Post("http://{serverAddressAndPort}/update/counter/{Name}/{Value}")

	if err != nil {
		panic(err)
	}

	fmt.Printf("counter ok %d", memstats.PullCount+1)
}

func CronRequest(pi time.Duration, ri time.Duration) {
	startTimePoll := time.Now()
	startTimeReport := time.Now()

	for {
		time.Sleep(time.Second)
		finishTimePoll := time.Now()
		finishTimeReport := time.Now()

		if finishTimePoll.Sub(startTimePoll) >= pi {
			memstats.CollectGaugeMetrics()
			startTimePoll = time.Now()
		}
		if finishTimeReport.Sub(startTimeReport) >= ri {
			SendMetrics2()
			startTimeReport = time.Now()
		}
	}
}

func main() {
	flag.StringVar(&serverAddr, "a", "localhost:8080", "server address and port, example 127.0.0.1:8080")
	flag.IntVar(&pollInterval, "p", defaultPollInterval, "poll interval, example 2 sec")
	flag.IntVar(&reportInterval, "r", defaultReportInterval, "report interval, example 10 sec")
	flag.Parse()

	CronRequest(time.Duration(pollInterval*NanoSecToSec), time.Duration(reportInterval*NanoSecToSec))
}

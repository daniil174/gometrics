package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/daniil174/gometrics/internal/memstats"
	"github.com/daniil174/gometrics/internal/storage"
	"github.com/go-resty/resty/v2"
)

const (
	RetryCount              = 5
	RetryMinWaitTimeSeconds = 5
	RetryMaxWaitTimeSeconds = 15
)

func compressData(data []byte) []byte {
	// Create a buffer to hold the compressed data.
	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)
	_, err := gzipWriter.Write(data)
	if err != nil {
		log.Fatalf("Failed to write to gzip writer: %v", err)
	}
	err = gzipWriter.Close()
	if err != nil {
		log.Fatalf("Failed to close gzip writer: %v", err)
	}
	return buf.Bytes()
}

func sendMetrics2(serverAddr string) {
	client := resty.New()

	for _, v := range memstats.CollectGaugeMetrics() {
		m := storage.Metrics{
			ID:    v.Name,
			MType: "gauge",
			Value: &v.Value,
		}

		req, rErr := json.Marshal(m)
		if rErr != nil {
			fmt.Println("Error occurred while making request:", rErr)
			return
		}

		_, err := client.
			SetRetryCount(RetryCount).
			SetRetryWaitTime(RetryMinWaitTimeSeconds*time.Second).
			SetRetryMaxWaitTime(RetryMaxWaitTimeSeconds*time.Second).
			R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Accept-Encoding", "gzip").
			SetHeader("Content-Encoding", "gzip").
			SetBody(compressData(req)).
			SetPathParams(map[string]string{
				"serverAddressAndPort": serverAddr,
			}).Post("http://{serverAddressAndPort}/update/")
		if err != nil {
			fmt.Println("Error occurred while making request:", err)
			return
		}
		println("gauge ok")
	}

	delta := memstats.PullCount + 1
	m := storage.Metrics{
		ID:    "PollCount",
		MType: "counter",
		Delta: &delta,
	}

	req, rErr := json.Marshal(m)
	if rErr != nil {
		fmt.Println("Error occurred while making request:", rErr)
		panic(rErr)
	}

	_, err := client.
		SetRetryCount(RetryCount).
		SetRetryWaitTime(RetryMinWaitTimeSeconds*time.Second).
		SetRetryMaxWaitTime(RetryMaxWaitTimeSeconds*time.Second).
		R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept-Encoding", "gzip").
		SetHeader("Content-Encoding", "gzip").
		SetBody(compressData(req)).
		SetPathParams(map[string]string{
			"serverAddressAndPort": serverAddr,
		}).Post("http://{serverAddressAndPort}/update/")

	if err != nil {
		panic(err)
	}

	fmt.Printf("counter ok %d", memstats.PullCount+1)
}

func CronRequest(pi time.Duration, ri time.Duration, servAddr string) {
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
			sendMetrics2(servAddr)
			startTimeReport = time.Now()
		}
	}
}

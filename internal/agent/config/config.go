package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v11"
)

const defaultPollInterval = 1
const defaultReportInterval = 3

var PollInterval int
var ReportInterval int
var ServerAddr string

type Config struct {
	Addr           string `env:"ADDRESS"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
}

func ConfigFromEnv() error {
	cfg, errConf := env.ParseAs[Config]()
	if errConf != nil {
		return errConf
	}
	fmt.Printf("ADDRESS=%s=", cfg.Addr)
	fmt.Printf("POLL_INTERVAL=%d=", cfg.PollInterval)
	fmt.Printf("REPORT_INTERVAL=%d=", cfg.ReportInterval)
	ServerAddr = cfg.Addr
	PollInterval = cfg.PollInterval
	ReportInterval = cfg.ReportInterval

	if ServerAddr == "" {
		flag.StringVar(&ServerAddr, "a", "localhost:8080", "server address and port, example 127.0.0.1:8080")
	}

	if PollInterval == 0 {
		flag.IntVar(&PollInterval, "p", defaultPollInterval, "poll interval, example 2 sec")
	}

	if ReportInterval == 0 {
		flag.IntVar(&ReportInterval, "r", defaultReportInterval, "report interval, example 10 sec")
	}

	flag.Parse()
	fmt.Printf("ServerAddr=%s=", ServerAddr)
	fmt.Printf("PollInterval=%d=", PollInterval)
	fmt.Printf("ReportInterval=%d= \n", ReportInterval)
	return nil
}

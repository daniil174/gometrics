package main

import (
	"github.com/daniil174/gometrics/internal/agent/agent"
	"github.com/daniil174/gometrics/internal/agent/config"
	"time"
)

const NanoSecToSec = 1000 * 1000 * 1000

func main() {
	_ = config.ConfigFromEnv()
	//ConfigFromEnv()

	agent.CronRequest(time.Duration(config.PollInterval*NanoSecToSec), time.Duration(config.ReportInterval*NanoSecToSec), config.ServerAddr)
}

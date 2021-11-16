package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

const (
	DefaultPollInterval = "30s"
	DefaultAppName      = "My Application"
	DefaultWorkloadName = "My Workload"
	NrMetricApi         = "https://metric-api.newrelic.com/metric/v1"
)

// To store configuration
type ConfigData struct {
	LicenseKey   string `json:"license_key"`
	PollInterval time.Duration
	Hostname     string
	Service      string
	Workload     string
}

func validateFile(file string) (err error) {
	var dirInfo os.FileInfo

	dirInfo, err = os.Stat(file)
	if os.IsNotExist(err) {
		err = fmt.Errorf("invalid, no such file [%s]", file)
		return
	} else if os.IsPermission(err) {
		err = fmt.Errorf("invalid, permission denied [%s]", file)
		return
	} else if err != nil {
		err = fmt.Errorf("invalid, %v [%s]", err, file)
		return
	} else if dirInfo.IsDir() {
		err = fmt.Errorf("invalid, filename is a directory [%s]", file)
	}
	return
}

func (data *ConfigData) initConfig() {
	var err error

	// Get license key
	data.LicenseKey = os.Getenv("NEW_RELIC_LICENSE_KEY")
	if len(data.LicenseKey) == 0 {
		log.Fatal("Error: could not locate env var NEW_RELIC_LICENSE_KEY")
	}
	data.Service = os.Getenv("NEW_RELIC_APP_NAME")
	if len(data.Service) == 0 {
		data.Service = DefaultAppName
	}
	data.Workload = os.Getenv("WORKLOAD_NAME")
	if len(data.Workload) == 0 {
		data.Workload = DefaultWorkloadName
	}

	// Get poll interval
	pollInterval := os.Getenv("POLL_INTERVAL")
	if len(pollInterval) == 0 {
		pollInterval = DefaultPollInterval
	}
	data.PollInterval, err = time.ParseDuration(pollInterval)
	if err != nil {
		log.Fatalf("Error: could not parse env var POLL_INTERVAL: %s, must me a duration (ex: 1h)", err)
	}

	// Get hostname
	data.Hostname, err = os.Hostname()
	if err != nil {
		log.Fatalf("Error: hostname of server %v", err)
	}

	log.Printf("Service: %s", data.Service)
	log.Printf("Workload: %s", data.Workload)
	log.Printf("Poll interval: %v", data.PollInterval)
}

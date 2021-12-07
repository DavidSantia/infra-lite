package main

import (
	"github.com/newrelic/infrastructure-agent/pkg/metrics/storage"
	"github.com/newrelic/infrastructure-agent/pkg/sample"
)

type DiskMonitor struct {
	samples         sample.EventBatch
	lastSamples     sample.EventBatch
	ioCountersStats map[string]storage.IOCountersStat
	lastDiskStats   map[string]storage.IOCountersStat
}

type DiskSample struct {
	UsedBytes               float64 `json:"diskUsedBytes"`
	UsedPercent             float64 `json:"diskUsedPercent"`
	FreeBytes               float64 `json:"diskFreeBytes"`
	FreePercent             float64 `json:"diskFreePercent"`
	TotalBytes              float64 `json:"diskTotalBytes"`
	UtilizationPercent      float64 `json:"diskUtilizationPercent"`
	ReadUtilizationPercent  float64 `json:"diskReadUtilizationPercent"`
	WriteUtilizationPercent float64 `json:"diskWriteUtilizationPercent"`
	ReadsPerSec             float64 `json:"diskReadsPerSecond"`
	WritesPerSec            float64 `json:"diskWritesPerSecond"`
}

func NewDiskMonitor() *DiskMonitor {
	return &DiskMonitor{}
}

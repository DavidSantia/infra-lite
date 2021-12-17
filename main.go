package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/newrelic/infrastructure-agent/pkg/sample"
)

type Metric map[string]interface{}

// To send JSON to NR Logs API
type Payload struct {
	Metrics []Metric `json:"metrics"`
}

// Make API request with error retry
func retryQuery(client *http.Client, method, url string, data []byte, headers []string) (b []byte) {
	var res *http.Response
	var err error
	var body io.Reader

	if len(data) > 0 {
		body = bytes.NewReader(data)
	}

	req, _ := http.NewRequest(method, url, body)
	for _, h := range headers {
		params := strings.Split(h, ":")
		req.Header.Set(params[0], params[1])
	}

	// up to 3 retries on API error
	for j := 1; j <= 3; j++ {
		res, err = client.Do(req)
		if err != nil {
			log.Println(err)
		} else if res.StatusCode == http.StatusOK || res.StatusCode == http.StatusAccepted {
			break
		} else {
			log.Printf("Retry %d: http status %d", j, res.StatusCode)
		}
	}
	if err == nil {
		b, err = ioutil.ReadAll(res.Body)
		if err != nil {
			log.Println(err)
			return
		}
		res.Body.Close()
	}
	return
}

func (data *ConfigData) makeMetric(name string, value float64) (metric Metric) {
	// Create metric API entry
	attributes := map[string]string{
		"workload": data.Workload,
		"service":  data.Service,
		"hostname": data.Hostname,
	}
	metric = Metric{
		"name":       data.Prefix + "." + name,
		"type":       "gauge",
		"value":      value,
		"timestamp":  data.SampleTime,
		"attributes": attributes,
	}
	return
}

func compressPayload(payload Payload) (b []byte) {
	// Marshall and compress JSON
	j, err2 := json.Marshal([]Payload{payload})
	if err2 != nil {
		log.Printf("Error formatting JSON for metrics api: %v", err2)
	}
	//log.Printf("Payload: %s", j)

	var gzBuf bytes.Buffer
	gz := gzip.NewWriter(&gzBuf)
	_, err2 = gz.Write(j)
	if err2 == nil {
		err2 = gz.Close()
	}
	if err2 != nil {
		log.Printf("Error compressing JSON for metrics api: %v", err2)
	}

	//log.Printf("Metric payload to post: length %d compressed %d", len(j), gzBuf.Len())
	b = gzBuf.Bytes()
	return
}

func main() {
	var err error
	var cpuSample *CPUSample
	var memSample *MemorySample
	var netSample, storageSample sample.EventBatch

	// Get configuration from env vars and/or newrelic.yml
	data := ConfigData{}
	data.initConfig()

	// Initialize monitors
	cpuMonitor := NewCPUMonitor()
	memoryMonitor := NewMemoryMonitor()
	networkMonitor := NewNetworkMonitor()
	storageMonitor := NewSampler(data.PollInterval)

	// Prime CPU and Disk monitor with first calls
	_, err = cpuMonitor.Sample()
	time.Sleep(time.Second)

	// Configure NR metrics API client
	client := &http.Client{}
	headers := []string{"Content-Type:application/json", "Content-Encoding:gzip", "Api-Key:" + data.LicenseKey}

	// Start poll loop
	for {
		startTime := time.Now()
		data.SampleTime = startTime.Unix()
		entries := make([]Metric, 0)

		// Fetch metrics
		cpuSample, err = cpuMonitor.Sample()
		if err != nil {
			log.Printf("Error: cpuMonitor %v", err)
		} else {
			entries = append(entries, data.makeMetric("CpuPercent", cpuSample.CPUPercent))
			entries = append(entries, data.makeMetric("CpuUserPercent", cpuSample.CPUUserPercent))
			entries = append(entries, data.makeMetric("CpuSystemPercent", cpuSample.CPUSystemPercent))
		}
		memSample, err = memoryMonitor.Sample()
		if err != nil {
			log.Printf("Error: memoryMonitor %v", err)
		} else {
			entries = append(entries, data.makeMetric("MemoryTotalBytes", memSample.MemoryTotal))
			entries = append(entries, data.makeMetric("MemoryFreeBytes", memSample.MemoryFree))
			entries = append(entries, data.makeMetric("MemoryUsedBytes", memSample.MemoryUsed))
			entries = append(entries, data.makeMetric("MemoryFreePercent", memSample.MemoryFreePercent))
			entries = append(entries, data.makeMetric("MemoryUsedPercent", memSample.MemoryUsedPercent))
			entries = append(entries, data.makeMetric("MemoryCachedBytes", memSample.MemoryCachedBytes))
			entries = append(entries, data.makeMetric("SwapTotalBytes", memSample.SwapTotal))
			entries = append(entries, data.makeMetric("SwapFreeBytes", memSample.SwapFree))
			entries = append(entries, data.makeMetric("SwapUsedBytes", memSample.SwapUsed))
		}

		netSample, err = networkMonitor.Sample()
		if err != nil {
			log.Printf("Error: networkMonitor %v", err)
		} else {
			for _, sample := range netSample {
				entries = append(entries, data.getNetworkMetric(sample, "ReceiveBytesPerSec"))
				entries = append(entries, data.getNetworkMetric(sample, "ReceiveErrorsPerSec"))
				entries = append(entries, data.getNetworkMetric(sample, "TransmitBytesPerSec"))
				entries = append(entries, data.getNetworkMetric(sample, "TransmitErrorsPerSec"))
			}
		}

		storageSample, err = storageMonitor.Sample()
		if err != nil {
			log.Printf("Error: storageMonitor %v", err)
		} else {
			for _, ss := range storageSample {
				entries = append(entries, data.getStorageMetric(ss, "UsedBytes"))
				entries = append(entries, data.getStorageMetric(ss, "UsedPercent"))
				entries = append(entries, data.getStorageMetric(ss, "FreeBytes"))
				entries = append(entries, data.getStorageMetric(ss, "FreePercent"))
				entries = append(entries, data.getStorageMetric(ss, "TotalBytes"))
				entries = append(entries, data.getStorageMetric(ss, "ReadBytesPerSec"))
				entries = append(entries, data.getStorageMetric(ss, "WriteBytesPerSec"))
				entries = append(entries, data.getStorageMetric(ss, "ReadWriteBytesPerSecond"))
			}
		}

		// Format for metrics API
		b := compressPayload(Payload{entries})

		// Post to API
		_ = retryQuery(client, "POST", NrMetricApi, b, headers)
		//log.Printf("Metrics api response %s", resp)

		remainder := data.PollInterval - time.Now().Sub(startTime)
		if remainder > 0 {
			//log.Printf("Sleeping %v", remainder)

			// Wait remainder of poll interval
			time.Sleep(remainder)
		}
	}
}

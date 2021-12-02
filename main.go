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
)

// To send JSON to NR Logs API
type Payload struct {
	Metrics []map[string]interface{} `json:"metrics"`
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

func (data *ConfigData) makeMetric(name string, value float64, ts time.Time) (metric map[string]interface{}) {
	// Create metric API entry
	attributes := map[string]string{
		"workload": data.Workload,
		"service":  data.Service,
		"hostname": data.Hostname,
	}
	metric = map[string]interface{}{
		"name":       name,
		"type":       "gauge",
		"value":      value,
		"timestamp":  ts.Unix(),
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

	// Get configuration from env vars and/or newrelic.yml
	data := ConfigData{}
	data.initConfig()

	// Initialize monitors
	cpuMonitor := NewCPUMonitor()
	memoryMonitor := NewMemoryMonitor()

	// Prime CPU monitor with first call
	_, err = cpuMonitor.Sample()
	time.Sleep(time.Second)

	// Configure NR metrics API client
	client := &http.Client{}
	headers := []string{"Content-Type:application/json", "Content-Encoding:gzip", "Api-Key:" + data.LicenseKey}

	// Start poll loop
	for {
		startTime := time.Now()
		entries := make([]map[string]interface{}, 0)

		// Fetch metrics
		cpuSample, err = cpuMonitor.Sample()
		if err != nil {
			log.Printf("Error: cpuMonitor %v", err)
			continue
		}
		entries = append(entries, data.makeMetric("cpuPercent", cpuSample.CPUPercent, startTime))
		entries = append(entries, data.makeMetric("cpuUserPercent", cpuSample.CPUUserPercent, startTime))
		entries = append(entries, data.makeMetric("cpuSystemPercent", cpuSample.CPUSystemPercent, startTime))
		memSample, err = memoryMonitor.Sample()
		if err != nil {
			log.Printf("Error: memoryMonitor %v", err)
			continue
		}
		entries = append(entries, data.makeMetric("memoryTotalBytes", memSample.MemoryTotal, startTime))
		entries = append(entries, data.makeMetric("memoryFreeBytes", memSample.MemoryFree, startTime))
		entries = append(entries, data.makeMetric("memoryUsedBytes", memSample.MemoryUsed, startTime))
		entries = append(entries, data.makeMetric("memoryFreePercent", memSample.MemoryFreePercent, startTime))
		entries = append(entries, data.makeMetric("memoryUsedPercent", memSample.MemoryUsedPercent, startTime))
		entries = append(entries, data.makeMetric("memoryCachedBytes", memSample.MemoryCachedBytes, startTime))
		entries = append(entries, data.makeMetric("swapTotalBytes", memSample.SwapTotal, startTime))
		entries = append(entries, data.makeMetric("swapFreeBytes", memSample.SwapFree, startTime))
		entries = append(entries, data.makeMetric("swapUsedBytes", memSample.SwapUsed, startTime))

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

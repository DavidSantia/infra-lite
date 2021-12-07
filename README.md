# infra-lite
Sends CPU and Memory metrics to NR metric api

## Overview
To use this utility, configure the following environment variables:

* NEW_RELIC_LICENSE_KEY
* NEW_RELIC_APP_NAME
* WORKLOAD_NAME
* POLL_INTERVAL
* METRIC_PREFIX

The license key env var is required.  The others have default values.

This utility will sample every 30s, pulling CPU and Memory metrics from the host or container.
You can adjust the POLL_INTERVAL as needed, to override the default 30s.

It will then send Metric data to the NR Metric API for each sample, as follows:

* container.CpuPercent
* container.CpuUserPercent
* container.CpuSystemPercent
* container.MemoryTotalBytes
* container.MemoryFreeBytes
* container.MemoryUsedBytes
* container.MemoryFreePercent
* container.MemoryUsedPercent
* container.MemoryCachedBytes
* container.SwapTotalBytes
* container.SwapFreeBytes
* container.SwapUsedBytes
* container.ReceiveBytesPerSec
* container.ReceiveErrorsPerSec
* container.TransmitBytesPerSec
* container.TransmitErrorsPerSec

Adjust the metric name prefix "container" if desired with the environment variable METRIC_PREFIX.
You can then query these in New Relic from the Metric namespace using NRQL.

## Build

Requires Go installed.  To build:
```sh
go get
go build
```


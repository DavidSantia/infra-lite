# infra-lite
Sends CPU, Memory, Network and Storage metrics to NR metric api

## Overview
To use this utility, configure the following environment variables:

* NEW_RELIC_LICENSE_KEY
* NEW_RELIC_APP_NAME
* NRIA_LOG_FILE
* NRIA_VERBOSE
* WORKLOAD_NAME
* POLL_INTERVAL
* METRIC_PREFIX

The `NEW_RELIC_LICENSE_KEY` environment variable is required.  The others have default values.

This utility will sample every 30s, pulling CPU, Memory, Network and Storage metrics from the host or container.
You can adjust `POLL_INTERVAL` as needed, to override the default 30s.

The agent will log startup information and errors to `./infra-lite.log` by default.
Use `NRIA_LOG_FILE` to override this filename. Set `NRIA_VERBOSE` to `1` to log any warnings.

It will then send the following Metric data to the NR Metric API:

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
* container.NetworkReceiveBytesPerSec
* container.NetworkReceiveErrorsPerSec
* container.NetworkTransmitBytesPerSec
* container.NetworkTransmitErrorsPerSec
* container.DiskUsedBytes
* container.DiskUsedPercent
* container.DiskFreeBytes
* container.DiskFreePercent
* container.DiskTotalBytes
* container.DiskReadBytesPerSec
* container.DiskWriteBytesPerSec
* container.DiskReadWriteBytesPerSecond

Adjust the metric name prefix "container" if desired with the environment variable `METRIC_PREFIX`

You can then query these meterics in NR1 from the Metric namespace using NRQL.

## Build

Requires Go installed.  To build:
```sh
go get
go build
```


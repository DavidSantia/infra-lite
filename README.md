# infra-lite
Sends CPU and Memory metrics to NR metric api

## Overview
To use this utility, configure the following environment variables:

* NEW_RELIC_LICENSE_KEY
* NEW_RELIC_APP_NAME
* WORKLOAD_NAME
* POLL_INTERVAL

The license key env var is required.  The others have default values.

This utility will sample every 30s, pulling CPU and Memory metrics from the host or container.
You can adjust the POLL_INTERVAL as needed, to override the default 30s.

It will then send Metric data to the NR Metric API for each sample, as follows:

* cpuPercent
* cpuUserPercent
* cpuSystemPercent
* memoryTotalBytes
* memoryFreeBytes
* memoryUsedBytes
* memoryFreePercent
* memoryUsedPercent
* memoryCachedBytes
* swapTotalBytes
* swapFreeBytes
* swapUsedBytes

You can then query these in New Relic from the Metric namespace using NRQL.

## Build

Requires Go installed.  To build:
```sh
go get
go build
```


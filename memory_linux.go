// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0
package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/mem"
)

// NewMemoryMonitor returns a memory monitor.
// It reports the free memory as the Available Memory, dependant on the current kernel
// or library implementations.
func NewMemoryMonitor() *MemoryMonitor {
	return &MemoryMonitor{vmHarvest: reclaimableAsUsed}
}

// Returns a formulation of the virtual memory that considers SReclaimable as Available, concretely:
// Total Memory: MemTotal
// Available Memory (kernels >= 3.14): MemAvailable
// Available Memory (kernels < 3.14): MemFree + Buffers + Cached
// Used Memory: Total Memory - Available Memory
func reclaimableAsUsed() (*mem.VirtualMemoryStat, error) {
	filename := "/proc/meminfo"
	return reclaimableAsUsedParseMemInfo(filename)
}

func reclaimableAsUsedParseMemInfo(filename string) (*mem.VirtualMemoryStat, error) {
	var i int
	var line string

	ret := &mem.VirtualMemoryStat{}

	f, err := os.Open(filename)
	if err != nil {
		log.Printf("Error parsing %s: %v", filename, err)
		return ret, nil
	}
	defer f.Close()
	r := bufio.NewReader(f)

	readFields := 0
	memAvailable := false
	for {
		line, err = r.ReadString('\n')
		if err == io.EOF || len(line) == 0 {
			break
		}
		fields := strings.Split(line, ":")
		if len(fields) != 2 {
			continue
		}
		key := strings.TrimSpace(fields[0])
		value := strings.TrimSpace(fields[1])
		value = strings.Replace(value, " kB", "", -1)

		t, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return ret, err
		}
		switch key {
		case "MemAvailable":
			ret.Available = t * 1024
			memAvailable = true
		case "MemTotal":
			ret.Total = t * 1024
			readFields++
		case "MemFree":
			ret.Free = t * 1024
			readFields++
		case "Buffers":
			ret.Buffers = t * 1024
			readFields++
		case "Cached":
			ret.Cached = t * 1024
			readFields++
		case "Shmem":
			ret.Shared = t * 1024
			readFields++
		case "Slab":
			ret.Slab = t * 1024
			readFields++
		case "SReclaimable":
			ret.SReclaimable = t * 1024
			readFields++
		}
		if readFields >= 7 && memAvailable { // stop reading the file when we have read all the fields we require
			break
		}
	}
	if !memAvailable {
		ret.Available = ret.Free + ret.Buffers + ret.Cached
	}
	ret.Used = ret.Total - ret.Available

	return ret, nil
}

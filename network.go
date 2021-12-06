package main

import (
	"github.com/newrelic/infrastructure-agent/pkg/metrics/network"
)

type NetworkEvent struct {
	InterfaceName        string
	HardwareAddress      string
	IpV4Address          string
	IpV6Address          string
	State                string
	ReceiveBytesPerSec   float64
	ReceiveErrorsPerSec  float64
	TransmitBytesPerSec  float64
	TransmitErrorsPerSec float64
}

func NewNetworkMonitor() network.NetworkSampler {
	return network.NetworkSampler{}
}

func (data *ConfigData) getNetworkMetric(ns *network.NetworkSample, name string) (metric Metric) {
	var value float64

	if name == "ReceiveBytesPerSec" && ns.ReceiveBytesPerSec != nil {
		value = *ns.ReceiveBytesPerSec
	} else if name == "ReceiveErrorsPerSec" && ns.ReceiveErrorsPerSec != nil {
		value = *ns.ReceiveErrorsPerSec
	} else if name == "TransmitBytesPerSec" && ns.TransmitBytesPerSec != nil {
		value = *ns.TransmitBytesPerSec
	} else if name == "TransmitErrorsPerSec" && ns.TransmitErrorsPerSec != nil {
		value = *ns.TransmitErrorsPerSec
	}
	metric = data.makeMetric(name, value)
	metric["attributes"].(map[string]string)["interfaceName"] = ns.InterfaceName
	metric["attributes"].(map[string]string)["hardwareAddress"] = ns.HardwareAddress
	metric["attributes"].(map[string]string)["ipV4Address"] = ns.IpV4Address
	metric["attributes"].(map[string]string)["ipV6Address"] = ns.IpV6Address
	metric["attributes"].(map[string]string)["state"] = ns.State
	return
}

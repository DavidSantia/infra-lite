package main

import (
	"github.com/newrelic/infrastructure-agent/pkg/metrics/network"
)

type NetworkSample *network.NetworkSample

func NewNetworkMonitor() network.NetworkSampler {
	return network.NetworkSampler{}
}

func (data *ConfigData) getNetworkMetric(sample interface{}, name string) (metric Metric) {
	var value float64
	ns := sample.(*network.NetworkSample)

	if name == "ReceiveBytesPerSec" && ns.ReceiveBytesPerSec != nil {
		value = *ns.ReceiveBytesPerSec
	} else if name == "ReceiveErrorsPerSec" && ns.ReceiveErrorsPerSec != nil {
		value = *ns.ReceiveErrorsPerSec
	} else if name == "TransmitBytesPerSec" && ns.TransmitBytesPerSec != nil {
		value = *ns.TransmitBytesPerSec
	} else if name == "TransmitErrorsPerSec" && ns.TransmitErrorsPerSec != nil {
		value = *ns.TransmitErrorsPerSec
	}
	metric = data.makeMetric("Network" + name, value)
	metric["attributes"].(map[string]string)["interfaceName"] = ns.InterfaceName
	metric["attributes"].(map[string]string)["hardwareAddress"] = ns.HardwareAddress
	metric["attributes"].(map[string]string)["ipV4Address"] = ns.IpV4Address
	metric["attributes"].(map[string]string)["ipV6Address"] = ns.IpV6Address
	metric["attributes"].(map[string]string)["state"] = ns.State
	return
}

package device

import (
	"github.com/cooper/screenmgr/measure"
	"strconv"
)

var hardwareOrder = [...]string{
	"ModelName",
	"ModelIdentifier",
	"CPUName",
	"CPUCount",
	"CPUCoreCount",
	"CPUFrequency",
	"L1Cache",
	"L2Cache",
	"L3Cache",
	"RAM",
	"BusFrequency",
}

var prettyHardware = map[string]string{
	"ModelName":       "Model Name",
	"ModelIdentifier": "Model Identifier",
	"CPUName":         "Processor Name",
	"CPUFrequency":    "Processor Speed",
	"CPUCount":        "Number of CPUs",
	"CPUCoreCount":    "Number of Cores",
	"L1Cache":         "L1 Cache (per CPU)",
	"L2Cache":         "L2 Cache (per CPU)",
	"L3Cache":         "L3 Cache (per CPU)",
	"RAM":             "Memory",
	"BusFrequency":    "Bus Frequency",
}

func mHzFromStringOrOne(str string) measure.Megahertz {
	n, err := strconv.ParseUint(str, 10, 0)
	if err != nil {
		return 1
	}
	return measure.Megahertz(n)
}

func (dev *Device) CPUMHz() measure.Megahertz {
	if freq, ok := dev.Info.Hardware["CPUFrequency"]; ok {
		return measure.MegahertzFromString(freq)
	}
	return 0
}

// CPU MHz multiplied to consider all CPUs/cores
func (dev *Device) CombinedCPUMHz() measure.Megahertz {
	theseMHz := dev.CPUMHz()

	// multiply for number of cpus
	if cpus, ok := dev.Info.Hardware["CPUCount"]; ok {
		theseMHz *= mHzFromStringOrOne(cpus)
	}

	// multiply for number of Cores
	if cores, ok := dev.Info.Hardware["CPUCoreCount"]; ok {
		theseMHz *= mHzFromStringOrOne(cores)
	}

	return theseMHz
}

func (dev *Device) RAMMBytes() measure.Megabytes {
	if ram, ok := dev.Info.Hardware["RAM"]; ok {
		return measure.MegabytesFromString(ram)
	}
	return 0
}

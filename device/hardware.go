package device

import (
	"github.com/cooper/screenmgr/measure"
	"strconv"
)

var hardwareOrder = [...]string{
	"ModelName",
	"ModelIdentifier",
	"CPUName",
	"CPUFrequency",
	"CPUCount",
	"CPUCoreCount",
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

type hardwarePair struct{ Key, Value string }

func (dev *Device) HardwareInOrder() (hardware []hardwarePair) {
	did := make(map[string]bool)
	doPair := func(key, val string) {
		if did[key] {
			return
		}
		hardware = append(hardware, hardwarePair{key, val})
		did[key] = true
	}

	// first, do the ordered ones
	for _, key := range hardwareOrder {
		if val, ok := dev.Info.Hardware[key]; ok {
			doPair(key, val)
		}
	}

	// then, do any extras
	// TODO: maybe alphabetize the leftovers?
	for key, val := range dev.Info.Hardware {
		doPair(key, val)
	}

	return
}

func (dev *Device) PrettyHardwareInOrder() (hardware []hardwarePair) {
	hardware = dev.HardwareInOrder()
	for i, pair := range hardware {
		if pretty, ok := prettyHardware[pair.Key]; ok {
			hardware[i].Key = pretty
		}
	}
	return
}

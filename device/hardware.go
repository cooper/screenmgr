package device

import (
	"github.com/cooper/screenmgr/measure"
	"strconv"
)

var hardwareOrder = [...]string{
	"ModelName",       // e.g. Power Mac G3
	"ModelIdentifier", // e.g. PowerMac1,1
	"CPUName",         // e.g. Intel Core i7 870
	"CPUFrequency",    // clock speed; e.g. 1.0 GHz
	"CPUCount",        // number of distinct processors
	"CPUCoreCount",    // number of cores PER PROCESSOR
	"L1Cache",         // e.g. 512 KB
	"L2Cache",         // e.g. 1 MB
	"L3Cache",         // e.g. 4 MB
	"RAM",             // system memory; e.g. 4 GB
	"BusFrequency",    // system bus speed; e.g. 100 MHz
}

var prettyHardware = map[string]interface{}{
	"ModelName":       "Model Name",
	"ModelIdentifier": "Model Identifier",
	"CPUName":         "Processor Name",
	"CPUFrequency":    "Processor Speed",
	"CPUCount":        "Number of CPUs",
	"CPUCoreCount":    "Cores (per CPU)",
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
		switch val := prettyHardware[pair.Key].(type) {

		case nil:
			break

		// a simple string is the pretty key name
		case string:
			hardware[i].Key = val

		// a function returns a pretty key and value
		case func(string, string) (string, string):
			hardware[i].Key, hardware[i].Value = val(pair.Key, pair.Value)

		// it can't be anything else
		default:
			panic("unknown type in prettyHardware map!")

		}
	}
	return
}

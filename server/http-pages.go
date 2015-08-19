package server

import (
	"fmt"
	"github.com/pivotal-golang/bytefmt"
	"strings"
)

type DevicePage struct {
	Devices []*Device
}

func (page *DevicePage) NumberOfDevices() int {
	return len(page.Devices)
}

func (page *DevicePage) RAMGBytes() string {
	var megabytes uint64 = 0
	for _, dev := range page.Devices {
		if ram, ok := dev.Info.Hardware["RAM"]; ok {
			megabytes += sizeStringToMegabytes(ram)
		}
	}
	return fmt.Sprintf("%.2f", float64(megabytes)/1024)
}

func (page *DevicePage) CPUGHz() string {
	return "0"
}

func sizeStringToMegabytes(size string) uint64 {
	size = strings.Replace(strings.TrimSuffix(size, "B"), " ", "", -1)
	mb, err := bytefmt.ToMegabytes(size)
	if err != nil {
		return 0
	}
	return mb
}

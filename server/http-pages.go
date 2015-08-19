package server

import (
	"fmt"
	"github.com/cooper/screenmgr/device"
	"github.com/pivotal-golang/bytefmt"
	"strings"
)

type DevicePage struct {
	Devices []*device.Device
}

func (page *DevicePage) NumberOfDevices() int {
	return len(page.Devices)
}

func (page *DevicePage) RAMGBytes() string {
	var megabytes uint64 = 0
	for _, dev := range page.Devices {
		fmt.Printf("checking if %s has it\n", dev.DeviceID)
		if ram, ok := dev.Info.Hardware["RAM"]; ok {
			fmt.Printf("...yes\n")
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
	fmt.Printf("so like: %s\n", size)
	mb, err := bytefmt.ToMegabytes(size)
	if err != nil {
		return 0
	}
	return mb
}

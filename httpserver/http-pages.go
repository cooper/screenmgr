package httpserver

import (
	"fmt"
	"github.com/cooper/screenmgr/device"
	"github.com/cooper/screenmgr/megabytes"
)

type DevicePage struct {
	Devices []*device.Device
}

func (page *DevicePage) NumberOfDevices() int {
	return len(page.Devices)
}

func (page *DevicePage) RAMGBytes() string {
	var mb float64 = 0
	for _, dev := range page.Devices {
		if ram, ok := dev.Info.Hardware["RAM"]; ok {
			mb += megabytes.MegabytesFromString(ram)
		}
	}
	return fmt.Sprintf("%.2f", megabytes.MegabytesToGigabytes(mb))
}

func (page *DevicePage) CPUGHz() string {
	return "0"
}

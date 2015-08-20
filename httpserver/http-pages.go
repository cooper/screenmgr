package httpserver

import (
	"github.com/cooper/screenmgr/device"
	"github.com/cooper/screenmgr/measure"
)

type DevicePage struct {
	Devices []*device.Device
}

func (page *DevicePage) NumberOfDevices() int {
	return len(page.Devices)
}

func (page *DevicePage) RAMGBytes() string {
	var mb measure.Megabytes
	for _, dev := range page.Devices {
		if ram, ok := dev.Info.Hardware["RAM"]; ok {
			mb += measure.MegabytesFromString(ram)
		}
	}
	return mb.ToGigabytes().String()
}

func (page *DevicePage) CPUGHz() string {
	return measure.GigahertzFromString("0 GHz").String()
}

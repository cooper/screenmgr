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
		mb += dev.RAMMBytes()
	}
	return mb.ToGigabytes().String()
}

func (page *DevicePage) CPUGHz() string {
	var mhz measure.Megahertz
	for _, dev := range page.Devices {
		mhz += dev.CombinedCPUMHz()
	}
	return mhz.ToGigahertz().String()
}

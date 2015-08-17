package server

type DevicePage struct {
	Devices []*Device
}

func (*DevicePage) NumberOfDevices() int {
    return len(devices)
}

func (*DevicePage) RAMGBytes() string {
    return "0"
}

func (*DevicePage) CPUGHz() string {
    return "0"
}

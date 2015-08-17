package server

type DevicePage struct {
	Devices []*Device
}

func (*DevicePage) NumberOfDevices() int {
    return len(devices)
}

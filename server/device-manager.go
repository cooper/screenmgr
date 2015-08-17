package server

import "io/ioutil"
import "fmt"

var devices = make(map[string]*Device)

// setup callbacks
type deviceSetupCallback func(dev *Device) error

var deviceSetupCallbacks []deviceSetupCallback

func addDeviceSetupCallback(cb deviceSetupCallback) {
	deviceSetupCallbacks = append(deviceSetupCallbacks, cb)
}

func findDevices() error {
	files, err := ioutil.ReadDir("devices")

	// error
	if err != nil {
		return err
	}

	for _, fileInfo := range files {
		fmt.Printf("found file: %+v\n", fileInfo)
		dev := deviceWithID(fileInfo.Name())

		// read the device info
		if err := dev.readInfo(); err != nil {
			reportError("Read JSON info", err)
			continue
		}

		// update device
		updateDevice(dev)

	}

	return nil
}

func setupDevices() error {
	for _, dev := range devices {
		err := dev.setup()
		if err != nil {
			return err
		}
	}
	return nil
}

// update the device in devices list with this device object,
// maybe preserving whatever we have already?
func updateDevice(dev *Device) {
	devices[dev.DeviceID] = dev
}

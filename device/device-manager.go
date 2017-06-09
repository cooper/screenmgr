package device

import (
	"io/ioutil"
	"log"
)

type DeviceSetupCallback func(dev *Device) error

var (
	devices              = make(map[string]*Device)
	deviceSetupCallbacks []DeviceSetupCallback
)

func AddDeviceSetupCallback(cb DeviceSetupCallback) {
	deviceSetupCallbacks = append(deviceSetupCallbacks, cb)
}

func FindDevices() error {
	files, err := ioutil.ReadDir("devices")

	// error
	if err != nil {
		return err
	}

	for _, fileInfo := range files {
		log.Printf("found file: %+v\n", fileInfo)
		dev := deviceWithID(fileInfo.Name())

		// read the device info
		if err := dev.ReadInfo(); err != nil {
			//reportError("Read JSON info", err)
			continue
		}

		// write it, just to beautify it
		dev.WriteInfo()

		// update device
		updateDevice(dev)

	}

	return nil
}

func SetupDevices() error {
	for _, dev := range devices {
		err := dev.setup()
		if err != nil {
			return err
		}
	}
	return nil
}

func AllDevices() map[string]*Device {
	return devices
}

// find an existing device by its ID
func GetDeviceByID(deviceID string) *Device {
	return devices[deviceID]
}

// update the device in devices list with this device object,
// maybe preserving whatever we have already?
func updateDevice(dev *Device) {
	devices[dev.DeviceID] = dev
}

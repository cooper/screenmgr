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

	count := 0
	for _, fileInfo := range files {
		count++
		dev := deviceWithID(fileInfo.Name())
		dev.Debug("initializing device")

		// read the device info
		if err := dev.ReadInfo(); err != nil {
			dev.Warn("%v", err)
			continue
		}

		// write it, just to beautify it
		dev.WriteInfo()

		// update device
		updateDevice(dev)
	}

	if count == 0 {
		log.Fatal("no devices configured")
	} else {
		log.Printf("looking for %d devices", count)
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

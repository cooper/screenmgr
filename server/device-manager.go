package server

import "io/ioutil"
import "fmt"

var devices = make(map[string]*device)

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

// update the device in devices list with this device object,
// maybe preserving whatever we have already?
func updateDevice(dev *device) {
	devices[dev.deviceID] = dev
}

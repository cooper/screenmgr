package server

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"os"
)

type deviceInfo struct {
	productLine string
	nickname    string
	addrString  string
}

type device struct {
	deviceID string
	info     deviceInfo
	addr     net.IP
}

// create a new device and a directory for it
func newDevice(deviceID string, info deviceInfo) (dev *device, err error) {

	// create device directory
	if err = createDeviceDirectory(deviceID); err != nil {
		return
	}

	// create the device object
	dev = deviceFromInfo(deviceID, info)

	// write the device data
	if err = dev.writeInfo(); err != nil {
		return
	}

	return
}

func deviceFromInfo(deviceID string, info deviceInfo) *device {
	return &device{}
	// TODO: parse IP here
}

// returns a directory for a device ID
func deviceDirectoryForID(deviceID string) string {
	return "devices/" + deviceID + ".device"
}

// create a new device directory
func createDeviceDirectory(deviceID string) error {
	return os.Mkdir(deviceDirectoryForID(deviceID), 0744)
}

// DEVICE METHODS

// get device directory
func (dev *device) getDirectory() string {
	return deviceDirectoryForID(dev.deviceID)
}

// get device info path
func (dev *device) getInfoPath() string {
	return dev.getDirectory() + "/info.json"
}

// write info to file
func (dev *device) writeInfo() error {
	json, err := json.Marshal(dev.info)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dev.getInfoPath(), json, 0744)
}

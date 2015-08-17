package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"os"
)

type deviceInfo struct {
	ProductLine string
	Nickname    string
	AddrString  string
	VNCEnabled  bool
	VNCPassword string
}

type device struct {
	deviceID string
	info     deviceInfo
}

// create a new device and a directory for it
func newDevice(deviceID string, info deviceInfo) (dev *device, err error) {

	// create device directory
	if err = createDeviceDirectory(deviceID); err != nil {
		return
	}

	// create the device object
	dev = deviceWithID(deviceID)
	dev.info = info

	// write the device data
	if err = dev.writeInfo(); err != nil {
		return
	}

	return
}

func deviceWithID(deviceID string) *device {
	return &device{deviceID: deviceID}
}

// returns a directory for a device ID
func deviceDirectoryForID(deviceID string) string {
	return "devices/" + deviceID
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

// get device file path
func (dev *device) getFilePath(fileName string) string {
	return dev.getDirectory() + "/" + fileName
}

// read info from file
func (dev *device) readInfo() error {
	data, err := ioutil.ReadFile(dev.getFilePath("info.json"))
	if err != nil {
		return err
	}
	if err = json.Unmarshal(data, &dev.info); err != nil {
		return err
	}
	return nil
}

// write info to file
func (dev *device) writeInfo() error {
	json, err := json.Marshal(dev.info)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dev.getFilePath("info.json"), json, 0744)
}

// get IP
func (dev *device) getIP() net.IP {
	return net.ParseIP(dev.info.AddrString)
}

// log warning
func (dev *device) warn(warning string) {
	log.Printf("[%s] %s\n", dev.deviceID, warning)
}

// setup device for loops and scrunch
func (dev *device) setup() error {
	for _, cb := range deviceSetupCallbacks {
		err := cb(dev)
		if err != nil {
			return nil
		}
	}
	return nil
}

package device

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
	Hardware    map[string]string
}

type Device struct {
	DeviceID       string
	Info           deviceInfo
	LastScreenshot string
	Online         bool
}

// create a new device and a directory for it
func NewDevice(deviceID string, info deviceInfo) (dev *Device, err error) {

	// create device directory
	if err = os.Mkdir(deviceDirectoryForID(deviceID), 0744); err != nil {
		return
	}

	// create the device object
	dev = deviceWithID(deviceID)
	dev.Info = info

	// create screenshot directory
	if err = os.Mkdir(dev.GetFilePath("screenshots"), 0744); err != nil {
		return
	}

	// write the device data
	if err = dev.WriteInfo(); err != nil {
		return
	}

	return
}

func deviceWithID(deviceID string) *Device {
	return &Device{DeviceID: deviceID}
}

// returns a directory for a device ID
func deviceDirectoryForID(deviceID string) string {
	return "devices/" + deviceID
}

// DEVICE METHODS

// get device directory
func (dev *Device) GetDirectory() string {
	return deviceDirectoryForID(dev.DeviceID)
}

// get device file path
func (dev *Device) GetFilePath(fileName string) string {
	return dev.GetDirectory() + "/" + fileName
}

// read info from file
func (dev *Device) ReadInfo() error {
	data, err := ioutil.ReadFile(dev.GetFilePath("info.json"))
	if err != nil {
		return err
	}
	if err = json.Unmarshal(data, &dev.Info); err != nil {
		return err
	}
	return nil
}

// write info to file
func (dev *Device) WriteInfo() error {
	json, err := json.Marshal(dev.Info)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dev.GetFilePath("info.json"), json, 0744)
}

// get IP
func (dev *Device) GetIP() net.IP {
	return net.ParseIP(dev.Info.AddrString)
}

// log warning
func (dev *Device) Warn(warning string) {
	log.Printf("[%s] %s\n", dev.DeviceID, warning)
}

// setup device for loops and scrunch
func (dev *Device) setup() error {
	for _, cb := range deviceSetupCallbacks {
		err := cb(dev)
		if err != nil {
			return nil
		}
	}
	return nil
}
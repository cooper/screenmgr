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
	Hardware    map[string]string
	Software    map[string]string

	VNC struct {
		Enabled  bool
		Password string
	}

	SSH struct {
		Enabled  bool
		Username string
		Password string
		UsesKey  bool
	}
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

// log
func (dev *Device) Log(f string, message ...interface{}) {
	log.Printf("["+dev.DeviceID+"] "+f+"\n", message...)
}

// log warning
func (dev *Device) Warn(f string, warning ...interface{}) {
	log.Printf("["+dev.DeviceID+"] "+f+"\n", warning...)
}

// log debug
func (dev *Device) Debug(f string, message ...interface{}) {
	log.Printf("["+dev.DeviceID+"] "+f+"\n", message...)
}

// find the last screenshot
func (dev *Device) GetLastScreenshot() string {

	// this is easy
	if dev.LastScreenshot != "" {
		return dev.LastScreenshot
	}

	// find screenshots
	files, err := ioutil.ReadDir(dev.GetFilePath("screenshots"))
	if err != nil {
		return ""
	}

	// find most recent
	var currentFile os.FileInfo
	var currentName string
	for _, file := range files {
		if currentFile == nil || file.ModTime().After(currentFile.ModTime()) {
			currentFile = file
			currentName = currentFile.Name()
		}
	}

	// cache this
	dev.LastScreenshot = currentName

	return currentName
}

// setup device for loops and such
func (dev *Device) setup() (err error) {
	for _, cb := range deviceSetupCallbacks {
		if err = cb(dev); err != nil {
			return
		}
	}
	return
}

package server

// #include "vncauth.h"
import "C"
import "os/exec"

// the VNC manager houses the devices in the VNC loop.
// not all devices in the master device list will
// necessarily be managed by the VNC manager.
type vncManager struct {
	devices []*device
}

var vnc = new(vncManager)

// start a loop for a device
func (vnc vncManager) startDeviceLoop(dev *device) {
	vnc.devices = append(vnc.devices, dev)

	// first check that it's enabled
	if !dev.info.VNCEnabled {
		dev.warn("Attempted to start VNC loop, but VNC is disabled")
		return
	}

	// check that there's a password
	if len(dev.info.VNCPassword) == 0 {
		dev.warn("Attempted to start VNC loop, but there's no password")
		return
	}

	// create a passwd file
	dir := dev.getFilePath("vncpasswd")
	C.vncEncryptAndStorePasswd(C.CString(dev.info.VNCPassword), C.CString(dir))

	// this method will loop so long as the device is
	// configured to run VNC.
	for dev.info.VNCEnabled {
		cmd := exec.Command("vncsnapshot",
			"-passwd", dir,
			"-fps", "5",
			"-count", "100",
            dev.info.AddrString,
            dev.getFilePath("screenshot.jpg"),
		)
		cmd.Start()
        break
	}

}

// add the VNC loop method to device setup
func init() {
	addDeviceSetupCallback(func(dev *device) error {
		vnc.startDeviceLoop(dev)
		return nil
	})
}

package server

import "os/exec"

// the VNC manager houses the devices in the VNC loop.
// not all devices in the master device list will
// necessarily be managed by the VNC manager.
type vncManager struct {
    devices []*device
}

// this is the initial run.
// other devices may be added later; this method
// merely starts the initial devices.
func (vnc vncManager) run() {
    for _, dev := range devices {
        vnc.startDeviceLoop(dev)
    }
}

func (vnc vncManager) startDeviceLoop(dev *device) {

    // TODO: generate password file

    // this method will loop so long as the device is
    // configured to run VNC.
    for dev.info.VNCEnabled {
        cmd := exec.Command("vncsnapshot", "")
        cmd.Start()
    }

}

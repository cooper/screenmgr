package vnc

// #include "vncauth.h"
import "C"

import (
	"bufio"
	"github.com/cooper/screenmgr/device"
	"os/exec"
	"regexp"
	"time"
)

// the VNC manager houses the devices in the VNC loop.
// not all devices in the master device list will
// necessarily be managed by the VNC manager.
type vncManager struct {
	devices []*device.Device
}

var vnc = new(vncManager)
var screenshotRegexp, _ = regexp.Compile(`screenshot\d+\.jpg`)

// start a loop for a device
func (vnc vncManager) startDeviceLoop(dev *device.Device) {
	vnc.devices = append(vnc.devices, dev)
	go vnc.deviceLoop(dev)
}

func (vnc vncManager) deviceLoop(dev *device.Device) {

	// first check that it's enabled
	if !dev.Info.VNCEnabled {
		dev.Warn("Attempted to start VNC loop, but VNC is disabled")
		return
	}

	// check that there's a password
	if len(dev.Info.VNCPassword) == 0 {
		dev.Warn("Attempted to start VNC loop, but there's no password")
		return
	}

	// create a passwd file
	passwd := dev.GetFilePath("vncpasswd")
	C.vncEncryptAndStorePasswd(C.CString(dev.Info.VNCPassword), C.CString(passwd))

	tryLater := func(errStr string) {
		dev.Online = false
		dev.Warn(errStr + "; waiting 10 seconds")
		time.Sleep(10 * time.Second)
	}

	// this method will loop so long as the device is
	// configured to run VNC.
VNCLoop:
	for dev.Info.VNCEnabled {
		cmd := exec.Command("vncsnapshot",
			"-passwd", passwd,
			"-fps", "5",
			"-count", "50",
			dev.Info.AddrString,
			dev.GetFilePath("screenshots/screenshot.jpg"),
		)

		// get STDERR and make a scanner
		stderr, err := cmd.StderrPipe()
		if err != nil {
			tryLater("failed to get vncsnapshot STDERR pipe")
			continue VNCLoop
		}
		scanner := bufio.NewScanner(stderr)

		// start the command
		if err := cmd.Start(); err != nil {
			tryLater("vncsnapshot failed to start")
			continue VNCLoop
		}

		// read from the scanner
		for scanner.Scan() {
			vnc.handleVNCSnapshotOutput(dev, scanner.Text())
		}

		// scanner error
		if err := scanner.Err(); err != nil {
			tryLater("Scanner terminated with error: " + err.Error())
			continue VNCLoop
		}

		// vncsnapshot exited with non-zero status
		if err := cmd.Wait(); err != nil {
			tryLater("vncsnapshot exited: " + err.Error())
			continue VNCLoop
		}

	}

}

func (vnc vncManager) handleVNCSnapshotOutput(dev *device.Device, line string) {
	found := screenshotRegexp.FindString(line)
	if len(found) == 0 {
		return
	}
	dev.Online = true
	dev.LastScreenshot = found
}

// add the VNC loop method to device setup
func init() {
	device.AddDeviceSetupCallback(func(dev *device.Device) error {
		vnc.startDeviceLoop(dev)
		return nil
	})
}

package server

// #include "vncauth.h"
import "C"

import (
	"bufio"
	"os/exec"
	"regexp"
	"time"
)

// the VNC manager houses the devices in the VNC loop.
// not all devices in the master device list will
// necessarily be managed by the VNC manager.
type vncManager struct {
	devices []*Device
}

var vnc = new(vncManager)
var screenshotRegexp, _ = regexp.Compile(`screenshot\d+\.jpg`)

// start a loop for a device
func (vnc vncManager) startDeviceLoop(dev *Device) {
	vnc.devices = append(vnc.devices, dev)
	go vnc.deviceLoop(dev)
}

func (vnc vncManager) deviceLoop(dev *Device) {

	// first check that it's enabled
	if !dev.Info.VNCEnabled {
		dev.warn("Attempted to start VNC loop, but VNC is disabled")
		return
	}

	// check that there's a password
	if len(dev.Info.VNCPassword) == 0 {
		dev.warn("Attempted to start VNC loop, but there's no password")
		return
	}

	// create a passwd file
	dir := dev.getFilePath("vncpasswd")
	C.vncEncryptAndStorePasswd(C.CString(dev.Info.VNCPassword), C.CString(dir))

	tryLater := func(errStr string) {
		dev.Online = false
		dev.warn(errStr + "; waiting 10 seconds")
		time.Sleep(10)
	}

	// this method will loop so long as the device is
	// configured to run VNC.
	for dev.Info.VNCEnabled {
		cmd := exec.Command("vncsnapshot",
			"-passwd", dir,
			"-fps", "5",
			"-count", "50",
			dev.Info.AddrString,
			dev.getFilePath("screenshots/screenshot.jpg"),
		)

		// get STDERR and make a scanner
		stderr, err := cmd.StderrPipe()
		if err != nil {
			tryLater("failed to get vncsnapshot STDERR pipe")
			continue
		}
		scanner := bufio.NewScanner(stderr)

		// start the command
		if err = cmd.Start(); err != nil {
			tryLater("vncsnapshot failed to start")
			continue
		}

		// read from the scanner
		for scanner.Scan() {
			vnc.handleVNCSnapshotOutput(dev, scanner.Text())
		}

		// scanner error
		if err := scanner.Err(); err != nil {
			tryLater("Scanner terminated with error: " + err.Error())
			continue
		}

		cmd.Wait()
	}

}

func (vnc vncManager) handleVNCSnapshotOutput(dev *Device, line string) {
	found := screenshotRegexp.FindString(line)
	if len(found) == 0 {
		return
	}
	dev.Online = true
	dev.LastScreenshot = found
}

// add the VNC loop method to device setup
func init() {

	addDeviceSetupCallback(func(dev *Device) error {
		vnc.startDeviceLoop(dev)
		return nil
	})
}

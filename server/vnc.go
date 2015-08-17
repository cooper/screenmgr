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
	devices []*device
}

var vnc = new(vncManager)
var screenshotRegexp, _ = regexp.Compile(`screenshot\d+\.jpg`)

// start a loop for a device
func (vnc vncManager) startDeviceLoop(dev *device) {
	vnc.devices = append(vnc.devices, dev)
	go vnc.deviceLoop(dev)
}

func (vnc vncManager) deviceLoop(dev *device) {

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

	tryLater := func(errStr string) {
		dev.warn(errStr + "; waiting 10 seconds")
		time.Sleep(10)
	}

	// this method will loop so long as the device is
	// configured to run VNC.
	for dev.info.VNCEnabled {
		cmd := exec.Command("vncsnapshot",
			"-passwd", dir,
			"-fps", "5",
			"-count", "50",
			dev.info.AddrString,
			dev.getFilePath("screenshot.jpg"),
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

	}

}

func (vnc vncManager) handleVNCSnapshotOutput(dev *device, line string) {
	found := screenshotRegexp.FindString(line)
	if len(found) == 0 {
		return
	}
	dev.lastScreenshot = found
}

// add the VNC loop method to device setup
func init() {

	addDeviceSetupCallback(func(dev *device) error {
		vnc.startDeviceLoop(dev)
		return nil
	})
}

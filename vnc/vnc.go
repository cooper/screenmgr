package vnc

// #include "vncauth.h"
import "C"

import (
	"bufio"
	"fmt"
	"github.com/cooper/screenmgr/device"
	"os/exec"
	"regexp"
	"time"
)

var devices []*device.Device
var screenshotRegexp, _ = regexp.Compile(`screenshot\d+\.jpg`)

// start a loop for a device
func startDeviceLoop(dev *device.Device) {
	devices = append(devices, dev)
	go deviceLoop(dev)
}

func deviceLoop(dev *device.Device) {

	// first check that it's enabled
	if !dev.Info.VNC.Enabled {
		dev.Warn("Attempted to start VNC loop, but VNC is disabled")
		return
	}

	// check that there's a password
	if len(dev.Info.VNC.Password) == 0 {
		dev.Warn("Attempted to start VNC loop, but there's no password")
		return
	}

	// create a passwd file
	passwd := dev.GetFilePath("vncpasswd")
	C.vncEncryptAndStorePasswd(C.CString(dev.Info.VNC.Password), C.CString(passwd))

	tryLater := func(errStr string) {
		dev.Warn("vnc: " + errStr + "; waiting 10 seconds")
		time.Sleep(10 * time.Second)
	}

	// this method will loop so long as the device is
	// configured to run VNC
	started := false

VNCLoop:
	for dev.Info.VNC.Enabled {

		// not online
		if !dev.Online {
			time.Sleep(10 * time.Second)
			continue VNCLoop
		}

		if !started {
			dev.Debug("starting VNC loop")
			started = true
		}

		// determine port
		port := dev.Info.VNC.Port
		if port == 0 {
			port = 5900
		}

		cmd := exec.Command("vncsnapshot",
			"-passwd", passwd,
			"-fps", "5",
			// "-encodings", "raw",
			"-count", "50",
			fmt.Sprintf("%s:%d", dev.Info.AddrString, port-5900),
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
			handleVNCSnapshotOutput(dev, scanner.Text())
		}

		// scanner error
		if err := scanner.Err(); err != nil {
			tryLater("scanner terminated with error: " + err.Error())
			continue VNCLoop
		}

		// vncsnapshot exited with non-zero status
		if err := cmd.Wait(); err != nil {
			tryLater("vncsnapshot exited: " + err.Error())
			continue VNCLoop
		}

	}

}

func handleVNCSnapshotOutput(dev *device.Device, line string) {
	dev.Debug("vncsnapshot: %s", line)
	found := screenshotRegexp.FindString(line)
	if len(found) == 0 {
		return
	}
	dev.LastScreenshot = found
	dev.Debug("updated screenshot: %s", found)
}

// add the VNC loop method to device setup
func init() {
	device.AddDeviceSetupCallback(func(dev *device.Device) error {
		startDeviceLoop(dev)
		return nil
	})
}

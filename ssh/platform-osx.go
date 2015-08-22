package ssh

import (
	"github.com/DHowett/go-plist"
	"github.com/cooper/screenmgr/device"
	"golang.org/x/crypto/ssh"
)

const (
	osxCommandVersion      = "sw_vers -productVersion"
	osxCommandHardwareInfo = "system_profiler -xml -detailLevel mini SPHardwareDataType"
)

// EXAMPLE OUTPUT
//
// Hardware Overview:
//
//   Model Name: Power Mac G5
//   Model Identifier: PowerMac7,2
//   Processor Name: PowerPC 970  (2.2)
//   Processor Speed: 2 GHz
//   Number Of CPUs: 2
//   L2 Cache (per CPU): 512 KB
//   Memory: 4 GB
//   Bus Speed: 1 GHz
//   Boot ROM Version: 5.1.4f0

func init() {
	initializers["osx"] = osxInitialize
}

func osxInitialize(dev *device.Device, sess *ssh.Session) error {
	data, err := sess.CombinedOutput(osxCommandHardwareInfo)
	if err != nil {
		dev.Warn("OS X system profiler command failed")
		return err
	}

	// unmarshal the plist
	var i interface{}
	if _, err := plist.Unmarshal(data, &i); err != nil {
		dev.Warn("failed to unmarshal system profiler plist")
		return err
	}

	return nil
}

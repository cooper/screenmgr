package ssh

import (
	"errors"
	"github.com/DHowett/go-plist"
	"github.com/cooper/screenmgr/device"
	"strconv"
)

const (
	osxCmdVersion      = `sw_vers -productVersion`
	osxCmdHardwareInfo = `system_profiler -xml -detailLevel mini SPHardwareDataType`
	osxCmdSerialNumber = `ioreg -c IOPlatformExpertDevice -d 2 | awk -F\" '/IOPlatformSerialNumber/{print $(NF-1)}'`
)

var osxVersionMap = [...]string{
	"", "", "", "",
	"Kodiak",
	"Puma",          // Darwin 5
	"Jaguar",        // Darwin 6
	"Panther",       // Darwin 7
	"Tiger",         // Darwin 8
	"Leopard",       // Darwin 9
	"Snow Leopard",  // Darwin 10
	"Lion",          // Darwin 11
	"Mountain Lion", // Darwin 12
	"Mavericks",     // Darwin 13
	"Yosemite",      // Darwin 14
	"El Capitan",    // Darwin 15
}

func init() {
	initializers["osx"] = osxInitialize
}

func osxInitialize(s sshClient) error {

	// run unix-like commands
	unixInitialize(s)

	// get OS X version information
	osxFindVersion(s)

	// extract the hardware info
	if err := osxFindHardware(s); err != nil {
		return err
	}

	return nil
}

var osxSystemProfilerMap = map[string]string{
	"physical_memory":         "RAM",
	"bus_speed":               "BusFrequency",
	"machine_model":           "ModelIdentifier",
	"machine_name":            "ModelName",
	"current_processor_speed": "CPUFrequency",
	"number_cpus":             "CPUCount",
	"number_processors":       "CPUCount",
	"cpu_type":                "CPUName",

	// "l2_cache_size":
	// "l2_cache_core":
	// "l3_cache":
	// "boot_rom_version":
	// "SMC_version_system":
}

func osxFindVersion(s sshClient) {
	dev := s.dev

	// get OS X numerical version
	dev.Info.Software["OSVersion"] = s.output(osxCmdVersion)

	// get OS X name; e.g. "Leopard"

	// determine the major version of Darwin
	darwinStr := ""
	for _, c := range dev.Info.Software["KernelVersion"] {
		if c == '.' {
			break
		}
		darwinStr += string(c)
	}

	// map Darwin version to OS X release
	if darwin, err := strconv.ParseUint(darwinStr, 10, 32); err != nil {
		dev.Warn("failed to parse Darwin version '%s': %s", darwinStr, err)
	} else {
		fullName := "OS X "
		if name := osxVersionMap[darwin]; name != "" {
			fullName += name
		} else {
			fullName += dev.Info.Software["OSVersion"]
		}
		dev.Info.Software["OSName"] = fullName
	}

}

func osxFindHardware(s sshClient) error {
	dev := s.dev

	// run system profiler
	data := s.outputBytes(osxCmdHardwareInfo)
	if data == nil {
		return errors.New("system profiler failed")
	}

	// unmarshal the plist
	var i interface{}
	if _, err := plist.Unmarshal(data, &i); err != nil {
		dev.Warn("failed to unmarshal system profiler plist")
		return err
	}

	// get the main array
	array, ok := i.([]interface{})
	if !ok {
		dev.Warn("expected a plist array")
		return errors.New("expected a plist array")
	}

	// get the main dict
	dict, ok := array[0].(map[string]interface{})
	if !ok {
		dev.Warn("expected a plist dictionary")
		return errors.New("expected a plist dictionary")
	}

	// get the results array
	array, ok = dict["_items"].([]interface{})
	if !ok {
		dev.Warn("expected a plist array")
		return errors.New("expected a plist array")
	}

	// finally, get the results dict
	dict, ok = array[0].(map[string]interface{})
	if !ok {
		dev.Warn("expected a plist dictionary")
		return errors.New("expected a plist dictionary")
	}

	dev.Debug("extracted from system profiler plist: %v", dict)
	osxHandleHardwareMap(dev, dict)

	return nil
}

func osxHandleHardwareMap(dev *device.Device, hw map[string]interface{}) {
	for key, val := range hw {
		str, ok := val.(string)
		if !ok || key == "_name" {
			continue
		}

		// convert system profiler key to screenmgr key
		if newKey, exists := osxSystemProfilerMap[key]; exists {
			dev.Info.Hardware[newKey] = str
			dev.Debug("%s = %s", newKey, str)
		} else {
			dev.Warn("ignoring system profiler key '%s'", key)
		}

	}
}

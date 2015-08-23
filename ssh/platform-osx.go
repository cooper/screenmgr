package ssh

import (
	"errors"
	"github.com/DHowett/go-plist"
	"github.com/cooper/screenmgr/device"
)

const (
	osxCommandVersion      = "sw_vers -productVersion"
	osxCommandHardwareInfo = "system_profiler -xml -detailLevel mini SPHardwareDataType"
)

func init() {
	initializers["osx"] = osxInitialize
}

func osxInitialize(s sshClient) error {
	dev := s.dev

	// run unix-like commands
	unixInitialize(s)

	// get OS X version
	dev.Info.Software["OSVersion"] = s.output(osxCommandVersion)

	// extract the hardware info
	data := s.outputBytes(osxCommandHardwareInfo)
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

func osxHandleHardwareMap(dev *device.Device, hw map[string]interface{}) {
	for key, val := range hw {
		str, ok := val.(string)
		if !ok || key == "_name" {
			continue
		}

		// convert to screenmgr key
		if newKey, exists := osxSystemProfilerMap[key]; exists {
			dev.Info.Hardware[newKey] = str
			dev.Debug("%s = %s", newKey, str)
		} else {
			dev.Warn("ignoring system profiler key '%s'", key)
		}

	}
}

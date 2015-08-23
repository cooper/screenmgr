package device

var softwareOrder = [...]string{
	"OSFamily",  // one of macos, osx, linux, bsd, windows
	"OSName",    // e.g. "OS X Leopard" or "Windows 8" or "Ubuntu"
	"OSVersion", // e.g. 10 for Windows 10, 10.4.3 for OS X Tiger
	"Kernel",
	"KernelVersion",
}

var prettySoftware = map[string]interface{}{
	"OSFamily": func(key, value string) (newKey, newValue string) {
		newKey = "OS Family"
		switch value {
		case "windows":
			newValue = "Windows"
		case "osx":
			newValue = "OS X"
		case "linux":
			newValue = "Linux"
		case "bsd":
			newValue = "BSD"
		case "macos":
			newValue = "Mac OS"
		default:
			newValue = value
		}
		return
	},
	"OSName":        "OS Name",
	"OSVersion":     "OS Version",
	"Kernel":        "Kernel",
	"KernelVersion": "Kernel Version",
}

type softwarePair struct{ Key, Value string }

func (dev *Device) SoftwareInOrder() (software []softwarePair) {
	did := make(map[string]bool)
	doPair := func(key, val string) {
		if did[key] {
			return
		}
		software = append(software, softwarePair{key, val})
		did[key] = true
	}

	// first, do the ordered ones
	for _, key := range softwareOrder {
		if val, ok := dev.Info.Software[key]; ok {
			doPair(key, val)
		}
	}

	// then, do any extras
	// TODO: maybe alphabetize the leftovers?
	for key, val := range dev.Info.Software {
		doPair(key, val)
	}

	return
}

func (dev *Device) PrettySoftwareInOrder() (software []softwarePair) {
	software = dev.SoftwareInOrder()
	for i, pair := range software {
		switch val := prettySoftware[pair.Key].(type) {

		case nil:
			break

		// a simple string is the pretty key name
		case string:
			software[i].Key = val

		// a function returns a pretty key and value
		case func(string, string) (string, string):
			software[i].Key, software[i].Value = val(pair.Key, pair.Value)

		// it can't be anything else
		default:
			panic("unknown type in prettySoftware map!")

		}
	}
	return
}

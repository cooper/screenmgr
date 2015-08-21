package device

import "sort"

type DeviceCollection []*Device

func (devs DeviceCollection) Len() int {
	return len(devs)
}

func (devs DeviceCollection) Less(i, j int) bool {
	dev1 := devs[i]
	dev2 := devs[j]

	// if only one is online, it comes first
	if dev1.Online && !dev2.Online {
		return true
	}
	if dev2.Online && !dev1.Online {
		return false
	}

	// otherwise, we fall back to alphabetical
	return sort.StringSlice{dev1.DeviceID, dev2.DeviceID}.Less(0, 1)

}

func (devs DeviceCollection) Swap(i, j int) {
	a := devs[i]
	b := devs[j]
	devs[j] = a
	devs[i] = b
}

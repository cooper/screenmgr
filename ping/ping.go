package ping

import (
	"github.com/cooper/screenmgr/device"
	"github.com/tatsushid/go-fastping"
	"net"
	"time"
)

var devices []*device.Device

func startDeviceLoop(dev *device.Device) error {
	devices = append(devices, dev)
	go deviceLoop(dev)
	return nil
}

func deviceLoop(dev *device.Device) {
	p := fastping.NewPinger()
	p.MaxRTT = 5 * time.Second
	p.Network("udp")

	// resolve IP
	ra, err := net.ResolveIPAddr("ip4:icmp", dev.Info.AddrString)
	if err != nil {
		dev.Warn("ping couldn't resolve the IP address")
		return
	}
	p.AddIPAddr(ra)

	// on receive, update last receive time
	lastTime := time.Now()
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		if !dev.Online {
			dev.Warn("now online")
			dev.Online = true
		}
		lastTime = time.Now()
	}

	// on idle, check if it's been a while
	p.OnIdle = func() {

		// it's been less than 10 seconds; no big deal
		if time.Since(lastTime) < 10*time.Second {
			return
		}

		if dev.Online {
			dev.Warn("now offline")
			dev.Online = false
		}
	}

	// do this continuously
	// TODO: if the device is removed or ping is disabled, stop the loop.
	p.RunLoop()

}

func init() {
	device.AddDeviceSetupCallback(startDeviceLoop)
}

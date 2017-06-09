package device

import "time"

func init() {
	AddDeviceSetupCallback(func(dev *Device) error {
		if !dev.Info.Autosave {
			return nil
		}
		go func() {
			for _ = range time.Tick(30 * time.Second) {
				dev.Debug("autosaving")
				dev.WriteInfo()
			}
		}()
		return nil
	})
}

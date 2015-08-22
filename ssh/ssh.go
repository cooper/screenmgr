package ssh

import (
	"github.com/cooper/screenmgr/device"
	"golang.org/x/crypto/ssh"
)

var devices []*device.Device
var initializers = make(map[string]func(dev *device.Device, sess *ssh.Session) error)

func startDeviceLoop(dev *device.Device) error {
	devices = append(devices, dev)
	go deviceLoop(dev)
	return nil
}

func deviceLoop(dev *device.Device) {
	if !dev.Info.SSH.Enabled {
		return
	}

	// create ssh config
	config := &ssh.ClientConfig{
		User: dev.Info.SSH.Username,
		Auth: authMethods(dev),
	}

	// dial
	// TODO: support other ports
	client, err := ssh.Dial("tcp", dev.Info.AddrString+":22", config)
	if err != nil {
		dev.Warn("ssh dial failed: %v", err)
		return
	}

	// create a session
	sess, err := client.NewSession()
	if err != nil {
		dev.Warn("ssh session init failed: %v", err)
		return
	}

	dev.Debug("SSH session established: %v", sess)
	initializers["osx"](dev, sess)
}

func authMethods(dev *device.Device) (methods []ssh.AuthMethod) {

	// TODO: keys
	if dev.Info.SSH.UsesKey {

	}

	if pw := dev.Info.SSH.Password; pw != "" {
		methods = append(methods, ssh.Password(pw))
	}

	return methods
}

func init() {
	device.AddDeviceSetupCallback(startDeviceLoop)
}

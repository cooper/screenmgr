package ssh

import (
	"github.com/cooper/screenmgr/device"
	"golang.org/x/crypto/ssh"
	"strings"
	"time"
)

// convenient for methods involving device + ssh client
type sshClient struct {
	dev    *device.Device
	client *ssh.Client
}

// the list of devices with ssh loop
var devices []*device.Device

// list of initializers for each OS family
// these are responsible for executing platform-specific commands
var initializers = make(map[string]func(sshClient) error)

func startDeviceLoop(dev *device.Device) error {
	devices = append(devices, dev)
	go deviceLoop(dev)
	return nil
}

func deviceLoop(dev *device.Device) {
	// TODO: actually loop on connection drop

	if !dev.Info.SSH.Enabled || dev.SSHRunning {
		return
	}

	dev.SSHRunning = true

	tryLater := func(errStr string) {
		dev.Warn("ssh: " + errStr + "; waiting 10 seconds")
		time.Sleep(10 * time.Second)
	}

	// create ssh config
	config := &ssh.ClientConfig{
		User: dev.Info.SSH.Username,
		Auth: authMethods(dev),
	}

SSHLoop:
	for dev.Info.SSH.Enabled {

		// not online
		if !dev.Online {
			time.Sleep(10 * time.Second)
			continue SSHLoop
		}

		// dial
		// TODO: support other ports
		addrStr := dev.Info.AddrString + ":22"
		dev.Debug("connecting ssh " + config.User + "@" + addrStr)
		client, err := ssh.Dial("tcp", addrStr, config)
		if err != nil {
			tryLater("ssh dial failed: " + err.Error())
			continue SSHLoop
		}

		dev.Log("ssh authenticated")

		// call initializer for this OS family
		family := dev.Info.Software["OSFamily"]
		if handler, exists := initializers[family]; exists {
			dev.Debug("initializing via ssh for OS family: %s", family)
			handler(sshClient{dev, client})
		} else {
			dev.Debug("no ssh handler for family " + family)
		}

		err = client.Wait()
		tryLater(err.Error())
	}

	dev.SSHRunning = false
}

// returns preferred authentication methods
func authMethods(dev *device.Device) (methods []ssh.AuthMethod) {

	// TODO: keys
	if dev.Info.SSH.UsesKey {

	}

	// password authentication
	if pw := dev.Info.SSH.Password; pw != "" {
		methods = append(methods, ssh.Password(pw))
	}

	return methods
}

// returns combined stdout + stderr
func (s sshClient) combinedOutputBytes(command string) []byte {
	sess, err := s.client.NewSession()
	if err != nil {
		s.dev.Warn("could not create ssh session: %v", err)
		return nil
	}
	data, err := sess.CombinedOutput(command)
	if err != nil {
		s.dev.Warn("command `%s` failed: %s", command, err)
		return nil
	}
	s.dev.Debug("`%s` = %s", command, data)
	return data
}

// returns combined stdout + stderr
func (s sshClient) combinedOutput(command string) string {
	return strings.TrimSpace(string(s.combinedOutputBytes(command)))
}

// returns stdout
func (s sshClient) outputBytes(command string) []byte {
	sess, err := s.client.NewSession()
	if err != nil {
		s.dev.Warn("could not create ssh session")
		return nil
	}
	data, err := sess.Output(command)
	if err != nil {
		s.dev.Warn("command `%s` failed: %s", command, err)
		return nil
	}
	s.dev.Debug("`%s` = %s", command, data)
	return data
}

// returns stdout
func (s sshClient) output(command string) string {
	return strings.TrimSpace(string(s.outputBytes(command)))
}

// returns exit code
func (s sshClient) command(command string) error {
	sess, err := s.client.NewSession()
	if err != nil {
		s.dev.Warn("could not create ssh session")
		return err
	}

	// run command
	err = sess.Run(command)

	// err != nil if exit code is non-zero
	s.dev.Debug("`%s`", command)
	if err != nil {
		s.dev.Warn("command `%s` failed: %s", command, err)
	}
	return err
}

func init() {
	device.AddDeviceSetupCallback(startDeviceLoop)
}

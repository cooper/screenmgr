package ssh

const (
	linuxCmdDistro  = `lsb_release -si`
	linuxCmdVersion = `lsb_release -sr`
)

func linuxInitialize(s sshClient) error {
	dev := s.dev
	unixInitialize(s)
	dev.Info.Software["OSName"] = s.output(linuxCmdDistro)
	dev.Info.Software["OSVersion"] = s.output(linuxCmdVersion)
	// TODO: if those are "", use other methods
	return nil
}

func init() {
	initializers["linux"] = linuxInitialize
}

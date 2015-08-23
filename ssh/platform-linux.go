package ssh

func linuxInitialize(s sshClient) error {
	dev := s.dev
	unixInitialize(s)
	dev.Info.Software["OSName"] = s.output("lsb_release -si")
	dev.Info.Software["OSVersion"] = s.output("lsb_release -sr")
	// TODO: if those are "", use other methods
	return nil
}

func init() {
	initializers["linux"] = linuxInitialize
}

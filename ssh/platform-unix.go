package ssh

func unixInitialize(s sshClient) error {
	dev := s.dev
	dev.Info.Software["Kernel"] = s.output("uname -s")
	dev.Info.Software["KernelVersion"] = s.output("uname -r")
	return nil
}

func init() {
	initializers["unix"] = unixInitialize
}

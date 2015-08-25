package ssh

const (
	unixCmdDirectory     = `mkdir -p $HOME/.screenmgr && cd $HOME/.screenmgr`
	unixCmdKernelName    = `uname -s`
	unixCmdKernelVersion = `uname -r`
)

func unixInitialize(s sshClient) error {
	dev := s.dev

	// make and enter a screenmgr working directory
	s.command(unixCmdDirectory)

	// get kernel info
	dev.Info.Software["Kernel"] = s.output(unixCmdKernelName)
	dev.Info.Software["KernelVersion"] = s.output(unixCmdKernelVersion)

	return nil
}

func init() {
	initializers["unix"] = unixInitialize
}

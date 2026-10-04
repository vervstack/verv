package server

type Options struct {
	UserName  string
	Password  string
	SshKeyUrl string
	// SysboxVersion is the sysbox release to install; empty means the newest numbered release.
	SysboxVersion string
	// IsDockerRestartAllowed confirms that installing sysbox may restart Docker and stop running containers.
	IsDockerRestartAllowed bool
}

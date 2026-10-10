package server

type Options struct {
	UserName  string
	Password  string
	SshKeyUrl string
	// SysboxVersion is the sysbox release to install; empty means the newest numbered release.
	SysboxVersion string
	// IsSysboxSkipped leaves sysbox uninstalled; the sysbox step reports itself as skipped.
	IsSysboxSkipped bool
	// IsDockerRestartAllowed confirms that installing sysbox may restart Docker and stop running containers.
	IsDockerRestartAllowed bool
}

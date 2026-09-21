package deploy

// platformAllowed reports whether deploy-velez may run on the given OS,
// honoring the --debug override for local development on non-Linux
// machines.
func platformAllowed(goos string, debug bool) bool {
	return goos == "linux" || debug
}

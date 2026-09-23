package menu

import (
	"os"
	"runtime"
)

// isHeadlessLinux reports whether verv is running on a Linux machine with no
// graphical session attached — the profile of a bare server someone is
// deploying a Velez node onto. DISPLAY/WAYLAND_DISPLAY are the standard X11
// and Wayland session markers; their absence on Linux is the conventional
// headless signal.
func isHeadlessLinux(goos, display, waylandDisplay string) bool {
	return goos == "linux" && display == "" && waylandDisplay == ""
}

// runningHeadlessLinux reads the real platform and environment.
func runningHeadlessLinux() bool {
	return isHeadlessLinux(runtime.GOOS, os.Getenv("DISPLAY"), os.Getenv("WAYLAND_DISPLAY"))
}

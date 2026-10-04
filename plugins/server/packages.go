package server

import (
	"context"
	"strings"
	"time"

	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/internal/cmd"
)

const (
	noninteractiveEnv = "DEBIAN_FRONTEND=noninteractive"

	// Any apt-get call can finish configuring a package left half-installed by an earlier
	// failure, and dpkg then asks about modified conffiles on a stdin that is closed.
	keepConfDefOption = "Dpkg::Options::=--force-confdef"
	keepConfOldOption = "Dpkg::Options::=--force-confold"

	// aptTimeout overrides cmd.Execute's 5s default, sized for quick admin
	// commands — package installs routinely run for minutes.
	aptTimeout = 10 * time.Minute

	dpkgInstalledStatus = "install ok installed"
)

var basePackages = []string{"curl", sudoGroup, "ca-certificates", "gnupg"}

func installBasePackages(_ context.Context, _ Options) (string, error) {
	err := aptGet("update")
	if err != nil {
		return "", rerrors.Wrap(err, "error updating package index")
	}

	installArgs := append([]string{"install", "-y"}, basePackages...)

	err = aptGet(installArgs...)
	if err != nil {
		return "", rerrors.Wrap(err, "error installing base packages")
	}

	return "Base packages installed", nil
}

func areBasePackagesInstalled(_ context.Context, _ Options) bool {
	for _, pkg := range basePackages {
		req := cmd.Request{Tool: "dpkg-query", Args: []string{"-W", "-f", "${Status}", pkg}}

		status, err := cmd.Execute(req)
		if err != nil || !isDpkgInstalled(status) {
			return false
		}
	}

	return true
}

func isDpkgInstalled(status string) bool {
	return strings.Contains(status, dpkgInstalledStatus)
}

func aptGet(args ...string) error {
	dpkgArgs := []string{"-o", keepConfDefOption, "-o", keepConfOldOption}

	req := cmd.Request{
		Tool:    "apt-get",
		Args:    append(dpkgArgs, args...),
		Env:     []string{noninteractiveEnv},
		Timeout: aptTimeout,
	}

	_, err := cmd.Execute(req)
	if err != nil {
		return rerrors.Wrap(err)
	}

	return nil
}

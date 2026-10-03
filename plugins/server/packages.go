package server

import (
	"context"
	"time"

	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/internal/cmd"
)

const (
	noninteractiveEnv = "DEBIAN_FRONTEND=noninteractive"

	// aptTimeout overrides cmd.Execute's 5s default, sized for quick admin
	// commands — package installs routinely run for minutes.
	aptTimeout = 10 * time.Minute
)

func installBasePackages(_ context.Context, _ Options) (string, error) {
	err := aptGet("update")
	if err != nil {
		return "", rerrors.Wrap(err, "error updating package index")
	}

	err = aptGet("install", "-y", "curl", "sudo", "ca-certificates", "gnupg")
	if err != nil {
		return "", rerrors.Wrap(err, "error installing base packages")
	}

	return "Base packages installed", nil
}

func aptGet(args ...string) error {
	req := cmd.Request{
		Tool:    "apt-get",
		Args:    args,
		Env:     []string{noninteractiveEnv},
		Timeout: aptTimeout,
	}

	_, err := cmd.Execute(req)
	if err != nil {
		return rerrors.Wrap(err)
	}

	return nil
}

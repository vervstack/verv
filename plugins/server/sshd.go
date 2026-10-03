package server

import (
	"context"
	"os"
	"strings"

	"go.redsock.ru/rerrors"
)

const (
	sshdConfigPath = "/etc/ssh/sshd_config"
	sshdConfigMode = 0o644
)

func enablePubkeyAuthInSshd(_ context.Context, _ Options) (string, error) {
	content, err := os.ReadFile(sshdConfigPath)
	if err != nil {
		return "", rerrors.Wrap(err, "error reading sshd config")
	}

	updated := enablePubkeyAuth(string(content))

	err = os.WriteFile(sshdConfigPath, []byte(updated), sshdConfigMode)
	if err != nil {
		return "", rerrors.Wrap(err, "error writing sshd config")
	}

	err = run("systemctl", "restart", "ssh")
	if err != nil {
		return "", rerrors.Wrap(err, "error restarting ssh")
	}

	return "SSH set up", nil
}

func enablePubkeyAuth(content string) string {
	return strings.ReplaceAll(content, "#PubkeyAuthentication yes", "PubkeyAuthentication yes")
}

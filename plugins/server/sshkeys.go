package server

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"go.redsock.ru/rerrors"
)

func installSshKeys(ctx context.Context, opts Options) (string, error) {
	userName := opts.UserName

	fetched, err := fetchText(ctx, opts.SshKeyUrl)
	if err != nil {
		return "", rerrors.Wrap(err, "error fetching ssh keys")
	}

	sshDir := filepath.Join(homeDir(userName), ".ssh")
	keysPath := filepath.Join(sshDir, "authorized_keys")

	err = os.MkdirAll(sshDir, 0o700)
	if err != nil {
		return "", rerrors.Wrap(err, "error creating .ssh directory")
	}

	existing, err := os.ReadFile(keysPath)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return "", rerrors.Wrap(err, "error reading authorized_keys")
	}

	merged := mergeAuthorizedKeys(string(existing), fetched)

	err = os.WriteFile(keysPath, []byte(merged), 0o600)
	if err != nil {
		return "", rerrors.Wrap(err, "error writing authorized_keys")
	}

	err = chownToUser(userName, sshDir, keysPath)
	if err != nil {
		return "", rerrors.Wrap(err, "error changing ownership")
	}

	return "SSH keys installed for " + userName, nil
}

func chownToUser(userName string, paths ...string) error {
	account, err := user.Lookup(userName)
	if err != nil {
		return rerrors.Wrap(err, "error looking up user")
	}

	uid, err := strconv.Atoi(account.Uid)
	if err != nil {
		return rerrors.Wrap(err, "error parsing uid")
	}

	gid, err := strconv.Atoi(account.Gid)
	if err != nil {
		return rerrors.Wrap(err, "error parsing gid")
	}

	for _, p := range paths {
		err = os.Chown(p, uid, gid)
		if err != nil {
			return rerrors.Wrap(err, "error changing owner of "+p)
		}
	}

	return nil
}

func mergeAuthorizedKeys(existing, fetched string) string {
	seen := make(map[string]struct{})

	for _, line := range strings.Split(existing, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			seen[trimmed] = struct{}{}
		}
	}

	var merged strings.Builder

	merged.WriteString(existing)

	if existing != "" && !strings.HasSuffix(existing, "\n") {
		merged.WriteString("\n")
	}

	for _, line := range strings.Split(fetched, "\n") {
		trimmed := strings.TrimSpace(line)

		_, isPresent := seen[trimmed]
		if trimmed == "" || isPresent {
			continue
		}

		seen[trimmed] = struct{}{}

		merged.WriteString(trimmed + "\n")
	}

	return merged.String()
}

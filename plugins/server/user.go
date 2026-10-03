package server

import (
	"context"
	"path/filepath"

	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/internal/cmd"
)

const homeRoot = "/home"

func homeDir(userName string) string {
	return filepath.Join(homeRoot, userName)
}

func UserExists(userName string) bool {
	err := run("id", "-u", userName)

	return err == nil
}

func ensureUser(_ context.Context, opts Options) (string, error) {
	doneMessage := "User " + opts.UserName + " already existed, password kept"

	if !UserExists(opts.UserName) {
		err := run("adduser", "--quiet", "--disabled-password", "--shell", "/bin/bash",
			"--home", homeDir(opts.UserName), "--gecos", opts.UserName, opts.UserName)
		if err != nil {
			return "", rerrors.Wrap(err, "error creating user")
		}

		doneMessage = "User " + opts.UserName + " created"
	}

	if opts.Password == "" {
		return doneMessage, nil
	}

	req := cmd.Request{
		Tool:  "chpasswd",
		Stdin: opts.UserName + ":" + opts.Password + "\n",
	}

	_, err := cmd.Execute(req)
	if err != nil {
		return "", rerrors.Wrap(err, "error setting password")
	}

	return doneMessage, nil
}

func addUserToGroups(_ context.Context, opts Options) (string, error) {
	userName := opts.UserName

	err := run("usermod", "-aG", "docker", userName)
	if err != nil {
		return "", rerrors.Wrap(err, "error adding user to docker group")
	}

	err = run("usermod", "-aG", "sudo", userName)
	if err != nil {
		return "", rerrors.Wrap(err, "error adding user to sudo group")
	}

	return userName + " added to docker and sudo groups", nil
}

package server

import (
	"context"
	"os"
	"strings"

	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/internal/cmd"
)

const (
	dockerGpgUrl   = "https://download.docker.com/linux/ubuntu/gpg"
	keyringsDir    = "/etc/apt/keyrings"
	dockerKeyPath  = "/etc/apt/keyrings/docker.gpg"
	dockerListPath = "/etc/apt/sources.list.d/docker.list"
	dockerRepoUrl  = "https://download.docker.com/linux/ubuntu"
)

func installDocker(ctx context.Context, _ Options) (string, error) {
	if cmd.IsInstalled("docker") {
		return "Docker already installed", nil
	}

	err := addDockerRepository(ctx)
	if err != nil {
		return "", rerrors.Wrap(err, "error adding docker repository")
	}

	err = aptGet("update")
	if err != nil {
		return "", rerrors.Wrap(err, "error updating package index")
	}

	args := []string{
		"install", "-y",
		"-o", "Dpkg::Options::=--force-confdef",
		"-o", "Dpkg::Options::=--force-confold",
		"docker-ce", "docker-ce-cli", "containerd.io", "docker-buildx-plugin", "docker-compose-plugin",
	}

	err = aptGet(args...)
	if err != nil {
		return "", rerrors.Wrap(err, "error installing docker packages")
	}

	return "Docker installed", nil
}

func addDockerRepository(ctx context.Context) error {
	err := os.MkdirAll(keyringsDir, 0o755)
	if err != nil {
		return rerrors.Wrap(err, "error creating keyrings directory")
	}

	key, err := fetchText(ctx, dockerGpgUrl)
	if err != nil {
		return rerrors.Wrap(err, "error fetching docker gpg key")
	}

	dearmor := cmd.Request{
		Tool:  "gpg",
		Args:  []string{"--batch", "--yes", "--dearmor", "-o", dockerKeyPath},
		Stdin: key,
	}

	_, err = cmd.Execute(dearmor)
	if err != nil {
		return rerrors.Wrap(err, "error dearmoring docker gpg key")
	}

	err = os.Chmod(dockerKeyPath, 0o644)
	if err != nil {
		return rerrors.Wrap(err, "error making docker gpg key readable")
	}

	printArch := cmd.Request{Tool: "dpkg", Args: []string{"--print-architecture"}}

	arch, err := cmd.Execute(printArch)
	if err != nil {
		return rerrors.Wrap(err, "error detecting architecture")
	}

	osRelease, err := os.ReadFile(osReleasePath)
	if err != nil {
		return rerrors.Wrap(err, "error reading os-release")
	}

	_, codename := parseOsRelease(string(osRelease))
	if codename == "" {
		return rerrors.Wrap(errUnknownCodename)
	}

	line := dockerRepoLine(strings.TrimSpace(arch), codename)

	err = os.WriteFile(dockerListPath, []byte(line), 0o644)
	if err != nil {
		return rerrors.Wrap(err, "error writing docker apt source")
	}

	return nil
}

func dockerRepoLine(arch, codename string) string {
	return "deb [arch=" + arch + " signed-by=" + dockerKeyPath + "] " + dockerRepoUrl + " " + codename + " stable\n"
}

func parseOsRelease(content string) (id, versionCodename string) {
	for _, line := range strings.Split(content, "\n") {
		key, value, found := strings.Cut(strings.TrimSpace(line), "=")
		if !found {
			continue
		}

		value = strings.Trim(value, `"'`)

		switch key {
		case "ID":
			id = value
		case "VERSION_CODENAME":
			versionCodename = value
		}
	}

	return id, versionCodename
}

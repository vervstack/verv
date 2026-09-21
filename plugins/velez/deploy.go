package velez

import (
	"fmt"
	"strconv"
	"strings"

	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/internal/cmd"
	"go.vervstack.ru/verv/internal/io"
)

const (
	containerName   = "velez"
	imageRepo       = "vervstack/velez"
	containerPort   = "53890"
	dockerSockMount = "/var/run/docker.sock:/var/run/docker.sock"
	diskMount       = "/dev/disk:/dev/disk"
	runMount        = "/run:/run"
)

// Deploy pulls the given vervstack/velez image tag, stops and removes any
// existing "velez" container, then starts the new one — the Go port of
// Velez's prod_init.sh.pattern deploy script, with two extra bind mounts
// (/dev/disk, /run) that its hardware-discovery library needs.
func Deploy(printer io.IO, version string, port int, keyPath string) error {
	image := imageRepo + ":" + version

	err := pullImage(printer, image)
	if err != nil {
		return rerrors.Wrap(err, "error pulling image")
	}

	err = stopExistingContainer(printer)
	if err != nil {
		return rerrors.Wrap(err, "error stopping existing container")
	}

	err = runContainer(printer, image, port, keyPath)
	if err != nil {
		return rerrors.Wrap(err, "error running container")
	}

	return nil
}

func pullImage(printer io.IO, image string) error {
	spinner := io.NewSpinner(printer)
	spinner.Start("Pulling " + image)

	_, err := cmd.Execute(cmd.Request{
		Tool: "docker",
		Args: []string{"pull", image},
	})
	if err != nil {
		spinner.Stop(false, "Pulling "+image+" — failed")

		return rerrors.Wrap(err)
	}

	spinner.Stop(true, "Pulled "+image)

	return nil
}

func stopExistingContainer(printer io.IO) error {
	spinner := io.NewSpinner(printer)
	spinner.Start("Checking for an existing " + containerName + " container")

	out, err := cmd.Execute(cmd.Request{
		Tool: "docker",
		Args: []string{"ps", "-aq", "-f", "name=" + containerName},
	})
	if err != nil {
		spinner.Stop(false, "Checking for an existing "+containerName+" container — failed")

		return rerrors.Wrap(err)
	}

	if strings.TrimSpace(out) == "" {
		spinner.Stop(true, "No existing "+containerName+" container found")

		return nil
	}

	spinner.Stop(true, "Found an existing "+containerName+" container")

	stopSpinner := io.NewSpinner(printer)
	stopSpinner.Start("Stopping and removing " + containerName)

	_, err = cmd.Execute(cmd.Request{
		Tool: "docker",
		Args: []string{"stop", containerName},
	})
	if err != nil {
		stopSpinner.Stop(false, "Stopping "+containerName+" — failed")

		return rerrors.Wrap(err)
	}

	_, err = cmd.Execute(cmd.Request{
		Tool: "docker",
		Args: []string{"rm", containerName},
	})
	if err != nil {
		stopSpinner.Stop(false, "Removing "+containerName+" — failed")

		return rerrors.Wrap(err)
	}

	stopSpinner.Stop(true, "Stopped and removed "+containerName)

	return nil
}

func runContainer(printer io.IO, image string, port int, keyPath string) error {
	spinner := io.NewSpinner(printer)
	spinner.Start("Starting " + containerName)

	portMapping := strconv.Itoa(port) + ":" + containerPort
	keyPathMount := keyPath + ":/tmp/velez"

	args := []string{
		"run",
		"-p", portMapping,
		"-d",
		"--name", containerName,
		"--restart=always",
		"-v", dockerSockMount,
		"-v", diskMount,
		"-v", runMount,
		"-v", keyPathMount,
		image,
	}

	_, err := cmd.Execute(cmd.Request{
		Tool: "docker",
		Args: args,
	})
	if err != nil {
		spinner.Stop(false, "Starting "+containerName+" — failed")

		return rerrors.Wrap(err)
	}

	spinner.Stop(true, fmt.Sprintf("Started %s on port %d", containerName, port))

	return nil
}

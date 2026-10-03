package velez

import (
	"encoding/json"
	"strconv"
	"strings"

	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/internal/cmd"
)

// Existing describes the settings of an already created "velez" container
// that an upgrade must carry over to the replacement container.
type Existing struct {
	Version string
	Port    int
	KeyPath string
}

type inspectedContainer struct {
	Config struct {
		Image string `json:"Image"`
	} `json:"Config"`
	HostConfig struct {
		PortBindings map[string][]struct {
			HostPort string `json:"HostPort"`
		} `json:"PortBindings"`
	} `json:"HostConfig"`
	Mounts []struct {
		Source      string `json:"Source"`
		Destination string `json:"Destination"`
	} `json:"Mounts"`
}

// FindExisting reports whether a "velez" container exists on this machine
// and, when it does, the version, host port and keys path it was started with.
// A machine without docker has no container, so that is not an error.
func FindExisting() (existing Existing, found bool, err error) {
	out, err := cmd.Execute(cmd.Request{
		Tool: "docker",
		Args: []string{"ps", "-aq", "-f", "name=^" + containerName + "$"},
	})
	if err != nil {
		return Existing{}, false, nil
	}

	if strings.TrimSpace(out) == "" {
		return Existing{}, false, nil
	}

	out, err = cmd.Execute(cmd.Request{
		Tool: "docker",
		Args: []string{"inspect", containerName},
	})
	if err != nil {
		return Existing{}, false, rerrors.Wrap(err, "error inspecting container")
	}

	existing, err = parseInspect(out)
	if err != nil {
		return Existing{}, false, rerrors.Wrap(err)
	}

	return existing, true, nil
}

func parseInspect(out string) (Existing, error) {
	var inspected []inspectedContainer

	err := json.Unmarshal([]byte(out), &inspected)
	if err != nil {
		return Existing{}, rerrors.Wrap(err, "error parsing container description")
	}

	if len(inspected) == 0 {
		return Existing{}, rerrors.Wrap(errEmptyInspect)
	}

	existing, err := toExisting(inspected[0])
	if err != nil {
		return Existing{}, rerrors.Wrap(err)
	}

	return existing, nil
}

func toExisting(c inspectedContainer) (Existing, error) {
	bindings := c.HostConfig.PortBindings[containerPort+"/tcp"]
	if len(bindings) == 0 {
		return Existing{}, rerrors.Wrap(errNoPortBinding)
	}

	port, err := strconv.Atoi(bindings[0].HostPort)
	if err != nil {
		return Existing{}, rerrors.Wrap(err, "error parsing host port")
	}

	keyPath := ""

	for _, m := range c.Mounts {
		if m.Destination == keysMountTarget {
			keyPath = m.Source
		}
	}

	if keyPath == "" {
		return Existing{}, rerrors.Wrap(errNoKeyMount)
	}

	version := ""

	tagIdx := strings.LastIndex(c.Config.Image, ":")
	if tagIdx >= 0 {
		version = c.Config.Image[tagIdx+1:]
	}

	return Existing{Version: version, Port: port, KeyPath: keyPath}, nil
}

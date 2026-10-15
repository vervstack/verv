package imagebuild

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"go.redsock.ru/rerrors"
)

const (
	dockerConfigEnv     = "DOCKER_CONFIG"
	dockerConfigDirName = ".docker"
	configFileName      = "config.json"
	temporaryDirPattern = "verv-docker-config-*"
	configFileMode      = 0o600
)

var credentialFields = []string{"credsStore", "credHelpers", "auths"}

type authEntry struct {
	Auth string `json:"auth"`
}

// A bare temporary DOCKER_CONFIG would lose docker's cli-plugins (buildx), contexts and buildx
// builders, which live in the real config dir, and `docker login` would persist the node secret
// there. So the temporary dir mirrors the real one via symlinks and only config.json is rewritten.
func prepareDockerConfig(credentials Credentials) (string, error) {
	realDir, err := realDockerConfigDir()
	if err != nil {
		return "", rerrors.Wrap(err, "error resolving docker config directory")
	}

	tempDir, err := os.MkdirTemp("", temporaryDirPattern)
	if err != nil {
		return "", rerrors.Wrap(err, "error creating temporary docker config directory")
	}

	err = populateDockerConfig(tempDir, realDir, credentials)
	if err != nil {
		removeDir(tempDir)
		return "", rerrors.Wrap(err, "error populating temporary docker config directory")
	}

	return tempDir, nil
}

func removeDir(dir string) {
	_ = os.RemoveAll(dir)
}

func realDockerConfigDir() (string, error) {
	fromEnv := os.Getenv(dockerConfigEnv)
	if fromEnv != "" {
		return fromEnv, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", rerrors.Wrap(err, "error resolving home directory")
	}

	return filepath.Join(home, dockerConfigDirName), nil
}

func populateDockerConfig(tempDir, realDir string, credentials Credentials) error {
	err := linkConfigEntries(tempDir, realDir)
	if err != nil {
		return rerrors.Wrap(err, "error linking docker config entries")
	}

	config, err := buildConfig(realDir, credentials)
	if err != nil {
		return rerrors.Wrap(err, "error building docker config")
	}

	err = os.WriteFile(filepath.Join(tempDir, configFileName), config, configFileMode)
	if err != nil {
		return rerrors.Wrap(err, "error writing docker config")
	}

	return nil
}

func linkConfigEntries(tempDir, realDir string) error {
	entries, err := os.ReadDir(realDir)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return rerrors.Wrap(err, "error reading docker config directory")
	}

	for _, entry := range entries {
		if entry.Name() == configFileName {
			continue
		}

		err = os.Symlink(filepath.Join(realDir, entry.Name()), filepath.Join(tempDir, entry.Name()))
		if err != nil {
			return rerrors.Wrap(err, "error linking entry")
		}
	}

	return nil
}

func buildConfig(realDir string, credentials Credentials) ([]byte, error) {
	fields, err := readConfigFields(filepath.Join(realDir, configFileName))
	if err != nil {
		return nil, rerrors.Wrap(err, "error reading docker config")
	}

	for _, field := range credentialFields {
		delete(fields, field)
	}

	entry := authEntry{Auth: encodeAuth(credentials)}
	auths := map[string]authEntry{registryHost(credentials.RegistryUrl): entry}

	rawAuths, err := json.Marshal(auths)
	if err != nil {
		return nil, rerrors.Wrap(err, "error marshaling auths")
	}
	fields["auths"] = rawAuths

	config, err := json.Marshal(fields)
	if err != nil {
		return nil, rerrors.Wrap(err, "error marshaling docker config")
	}

	return config, nil
}

func readConfigFields(path string) (map[string]json.RawMessage, error) {
	fields := make(map[string]json.RawMessage)

	raw, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return fields, nil
	}
	if err != nil {
		return nil, rerrors.Wrap(err, "error reading file")
	}

	err = json.Unmarshal(raw, &fields)
	if err != nil {
		return nil, rerrors.Wrap(err, "error parsing json")
	}

	return fields, nil
}

func encodeAuth(credentials Credentials) string {
	return base64.StdEncoding.EncodeToString([]byte(credentials.Username + ":" + credentials.Password))
}

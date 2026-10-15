package projectinfo

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"go.redsock.ru/rerrors"
	"gopkg.in/yaml.v3"
)

const vervonomiconFileName = "vervonomicon.yaml"

type vervonomicon struct {
	Service serviceSection `yaml:"service"`
	Name    string         `yaml:"name"`
}

type serviceSection struct {
	Name string `yaml:"name"`
}

func ServiceName(projectDir string) (string, error) {
	path := filepath.Join(projectDir, vervDirName, vervonomiconFileName)

	raw, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return "", rerrors.Wrap(ErrServiceNameUnknown)
	}
	if err != nil {
		return "", rerrors.Wrap(err, "error reading vervonomicon")
	}

	var content vervonomicon
	err = yaml.Unmarshal(raw, &content)
	if err != nil {
		return "", rerrors.Wrap(err, "error parsing vervonomicon")
	}

	rawName := content.Service.Name
	if rawName == "" {
		rawName = content.Name
	}

	name := normalizeServiceName(rawName)
	if name == "" {
		return "", rerrors.Wrap(ErrServiceNameUnknown)
	}

	return name, nil
}

func normalizeServiceName(raw string) string {
	trimmed := strings.TrimRight(strings.TrimSpace(raw), "/")
	lastSegment := trimmed[strings.LastIndex(trimmed, "/")+1:]

	return strings.Map(keepServiceNameRune, strings.ToLower(lastSegment))
}

func keepServiceNameRune(r rune) rune {
	isLetter := r >= 'a' && r <= 'z'
	isDigit := r >= '0' && r <= '9'
	isSeparator := r == '.' || r == '_' || r == '-'

	if isLetter || isDigit || isSeparator {
		return r
	}

	return -1
}

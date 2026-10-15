package projectinfo

import (
	"cmp"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"go.redsock.ru/rerrors"
)

const (
	vervDirName        = ".verv"
	configDirName      = "config"
	configExtension    = ".yaml"
	prodConfigName     = "config"
	prodEnvironment    = "prod"
	templateConfigName = "config_template"
)

var environmentPriority = map[string]int{
	"stage": 0,
	"local": 1,
	"prod":  2,
}

func Environments(projectDir string) ([]string, error) {
	found := make(map[string]struct{})

	configFiles, err := listEntryNames(filepath.Join(projectDir, configDirName), false)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing config files")
	}

	for _, fileName := range configFiles {
		if filepath.Ext(fileName) != configExtension {
			continue
		}

		name := strings.TrimSuffix(fileName, configExtension)
		if name == templateConfigName {
			continue
		}
		if name == prodConfigName {
			name = prodEnvironment
		}

		found[name] = struct{}{}
	}

	vervDirs, err := listEntryNames(filepath.Join(projectDir, vervDirName), true)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing verv directories")
	}

	for _, dirName := range vervDirs {
		if strings.HasPrefix(dirName, ".") {
			continue
		}

		found[dirName] = struct{}{}
	}

	environments := make([]string, 0, len(found))
	for name := range found {
		environments = append(environments, name)
	}
	slices.SortFunc(environments, compareEnvironments)

	return environments, nil
}

func listEntryNames(dir string, isDirWanted bool) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, rerrors.Wrap(err, "error reading directory")
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() == isDirWanted {
			names = append(names, entry.Name())
		}
	}

	return names, nil
}

func compareEnvironments(a, b string) int {
	byPriority := cmp.Compare(environmentPriorityOf(a), environmentPriorityOf(b))
	if byPriority != 0 {
		return byPriority
	}

	return cmp.Compare(a, b)
}

func environmentPriorityOf(name string) int {
	priority, ok := environmentPriority[name]
	if ok {
		return priority
	}

	return len(environmentPriority)
}

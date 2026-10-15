package nodes

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"go.redsock.ru/rerrors"
)

const (
	gitignoreFileName = ".gitignore"
	gitignoreEntry    = ".verv/nodes.local.yaml"
)

func ensureGitignored(projectDir string) error {
	path := filepath.Join(projectDir, gitignoreFileName)

	content, err := os.ReadFile(path)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return rerrors.Wrap(err, "error reading gitignore")
	}

	if hasLine(string(content), gitignoreEntry) {
		return nil
	}

	updated := string(content)
	if updated != "" && !strings.HasSuffix(updated, "\n") {
		updated += "\n"
	}
	updated += gitignoreEntry + "\n"

	err = os.WriteFile(path, []byte(updated), 0o644)
	if err != nil {
		return rerrors.Wrap(err, "error writing gitignore")
	}

	return nil
}

func hasLine(content, line string) bool {
	for _, candidate := range strings.Split(content, "\n") {
		if strings.TrimSpace(candidate) == line {
			return true
		}
	}

	return false
}

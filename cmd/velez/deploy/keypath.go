package deploy

import (
	"os"
	"path/filepath"
	"strings"

	"go.redsock.ru/rerrors"
)

// expandKeyPath resolves a leading "~" in path against the current user's
// home directory. docker's -v flag does not expand "~" itself, so this must
// happen before the path reaches the docker args. A path with no leading
// "~" is returned unchanged.
func expandKeyPath(path string) (string, error) {
	if path != "~" && !strings.HasPrefix(path, "~/") {
		return path, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", rerrors.Wrap(err, "error resolving home directory")
	}

	if path == "~" {
		return home, nil
	}

	return filepath.Join(home, strings.TrimPrefix(path, "~/")), nil
}

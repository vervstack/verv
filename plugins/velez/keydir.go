package velez

import (
	"errors"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"strconv"

	"go.redsock.ru/rerrors"
)

const (
	DefaultKeyPath   = "/opt/velez"
	keyDirGroup      = "verv"
	KeyDirSharedMode = 0o770
	keyDirPlainMode  = 0o755
	rootUid          = 0
	homeLinkName     = "velez"
)

func PrepareKeyDir(keyPath string) error {
	if keyPath != DefaultKeyPath {
		return nil
	}

	if os.Geteuid() != rootUid {
		info, err := os.Stat(keyPath)
		if err == nil && info.IsDir() {
			return nil
		}

		return rerrors.Wrap(errKeyDirNeedsRoot)
	}

	err := os.MkdirAll(keyPath, keyDirPlainMode)
	if err != nil {
		return rerrors.Wrap(err, "error creating key dir")
	}

	group, err := user.LookupGroup(keyDirGroup)
	if err != nil {
		err = os.Chmod(keyPath, keyDirPlainMode)
		if err != nil {
			return rerrors.Wrap(err, "error setting key dir mode")
		}

		return nil
	}

	gid, err := strconv.Atoi(group.Gid)
	if err != nil {
		return rerrors.Wrap(err, "error parsing group id")
	}

	err = os.Chown(keyPath, rootUid, gid)
	if err != nil {
		return rerrors.Wrap(err, "error changing key dir owner")
	}

	err = os.Chmod(keyPath, KeyDirSharedMode)
	if err != nil {
		return rerrors.Wrap(err, "error setting key dir mode")
	}

	return nil
}

func LinkInHome(keyPath string) error {
	if keyPath != DefaultKeyPath {
		return nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return rerrors.Wrap(err, "error resolving home dir")
	}

	link := filepath.Join(home, homeLinkName)

	_, err = os.Lstat(link)
	if err == nil {
		return nil
	}

	if !errors.Is(err, fs.ErrNotExist) {
		return rerrors.Wrap(err, "error checking home link")
	}

	err = os.Symlink(keyPath, link)
	if err != nil {
		return rerrors.Wrap(err, "error creating home link")
	}

	return nil
}

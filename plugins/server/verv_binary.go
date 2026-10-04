package server

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"io"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"

	"go.redsock.ru/rerrors"
)

const (
	vervGroup       = "verv"
	vervInstallDir  = "/usr/local/bin"
	vervInstallPath = vervInstallDir + "/verv"
	vervTempPattern = ".verv-install-*"
	vervBinaryMode  = 0o750
	rootUid         = 0
)

func isVervBinaryInstalled(_ context.Context, opts Options) bool {
	info, err := os.Stat(vervInstallPath)
	if err != nil || info.Mode().Perm() != vervBinaryMode {
		return false
	}

	group, err := user.LookupGroup(vervGroup)
	if err != nil {
		return false
	}

	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok || strconv.FormatUint(uint64(stat.Gid), 10) != group.Gid {
		return false
	}

	if !isUserInGroups(opts.UserName, vervGroup) {
		return false
	}

	return isInstalledBinaryCurrent()
}

func isInstalledBinaryCurrent() bool {
	isInstalled, err := isRunningFromInstallPath()
	if err != nil {
		return false
	}

	if isInstalled {
		return true
	}

	source, err := runningBinaryPath()
	if err != nil {
		return false
	}

	sourceSum, err := fileSha256(source)
	if err != nil {
		return false
	}

	installedSum, err := fileSha256(vervInstallPath)
	if err != nil {
		return false
	}

	return bytes.Equal(sourceSum, installedSum)
}

func fileSha256(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, rerrors.Wrap(err, "error opening "+path)
	}

	defer func() { _ = file.Close() }()

	hash := sha256.New()

	_, err = io.Copy(hash, file)
	if err != nil {
		return nil, rerrors.Wrap(err, "error hashing "+path)
	}

	return hash.Sum(nil), nil
}

func installVervBinary(_ context.Context, opts Options) (string, error) {
	err := run("groupadd", "-f", vervGroup)
	if err != nil {
		return "", rerrors.Wrap(err, "error creating group "+vervGroup)
	}

	err = run("usermod", "-aG", vervGroup, opts.UserName)
	if err != nil {
		return "", rerrors.Wrap(err, "error adding user to group "+vervGroup)
	}

	group, err := user.LookupGroup(vervGroup)
	if err != nil {
		return "", rerrors.Wrap(err, "error looking up group "+vervGroup)
	}

	gid, err := strconv.Atoi(group.Gid)
	if err != nil {
		return "", rerrors.Wrap(err, "error parsing gid of group "+vervGroup)
	}

	isInstalled, err := isRunningFromInstallPath()
	if err != nil {
		return "", rerrors.Wrap(err)
	}

	if !isInstalled {
		err = copyRunningBinary(gid)
		if err != nil {
			return "", rerrors.Wrap(err)
		}
	}

	err = restrictInstalledBinary(gid)
	if err != nil {
		return "", rerrors.Wrap(err)
	}

	return "verv installed to " + vervInstallPath + " (group " + vervGroup + ": root and " + opts.UserName +
		"); log in again for the group to apply", nil
}

func isRunningFromInstallPath() (bool, error) {
	source, err := runningBinaryPath()
	if err != nil {
		return false, rerrors.Wrap(err)
	}

	dest, err := filepath.EvalSymlinks(vervInstallPath)
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}

	if err != nil {
		return false, rerrors.Wrap(err, "error resolving "+vervInstallPath)
	}

	return source == dest, nil
}

func runningBinaryPath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", rerrors.Wrap(err, "error getting executable path")
	}

	resolved, err := filepath.EvalSymlinks(exe)
	if err != nil {
		return "", rerrors.Wrap(err, "error resolving executable symlinks")
	}

	return resolved, nil
}

// copyRunningBinary writes through a temp file and renames it into place: the rename is atomic and
// safe when the destination is a binary that is currently running.
func copyRunningBinary(gid int) error {
	source, err := runningBinaryPath()
	if err != nil {
		return rerrors.Wrap(err)
	}

	in, err := os.Open(source)
	if err != nil {
		return rerrors.Wrap(err, "error opening running binary")
	}

	defer func() { _ = in.Close() }()

	tmp, err := os.CreateTemp(vervInstallDir, vervTempPattern)
	if err != nil {
		return rerrors.Wrap(err, "error creating temp file in "+vervInstallDir)
	}

	tmpPath := tmp.Name()

	err = fillTempBinary(tmp, in, gid)
	if err != nil {
		_ = os.Remove(tmpPath)

		return rerrors.Wrap(err)
	}

	err = os.Rename(tmpPath, vervInstallPath)
	if err != nil {
		_ = os.Remove(tmpPath)

		return rerrors.Wrap(err, "error moving binary to "+vervInstallPath)
	}

	return nil
}

func fillTempBinary(tmp *os.File, in io.Reader, gid int) error {
	_, err := io.Copy(tmp, in)
	if err != nil {
		_ = tmp.Close()

		return rerrors.Wrap(err, "error copying binary")
	}

	err = tmp.Close()
	if err != nil {
		return rerrors.Wrap(err, "error closing temp file")
	}

	err = restrictFile(tmp.Name(), gid)
	if err != nil {
		return rerrors.Wrap(err)
	}

	return nil
}

func restrictInstalledBinary(gid int) error {
	err := restrictFile(vervInstallPath, gid)
	if err != nil {
		return rerrors.Wrap(err)
	}

	return nil
}

func restrictFile(path string, gid int) error {
	err := os.Chown(path, rootUid, gid)
	if err != nil {
		return rerrors.Wrap(err, "error changing owner of "+path)
	}

	err = os.Chmod(path, vervBinaryMode)
	if err != nil {
		return rerrors.Wrap(err, "error changing mode of "+path)
	}

	return nil
}

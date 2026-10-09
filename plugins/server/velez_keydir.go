package server

import (
	"context"
	"os"
	"os/user"
	"strconv"
	"syscall"

	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/plugins/velez"
)

func isVelezKeyDirSetUp(_ context.Context, _ Options) bool {
	info, err := os.Stat(velez.DefaultKeyPath)
	if err != nil || !info.IsDir() || info.Mode().Perm() != velez.KeyDirSharedMode {
		return false
	}

	group, err := user.LookupGroup(vervGroup)
	if err != nil {
		return false
	}

	stat, ok := info.Sys().(*syscall.Stat_t)

	return ok && strconv.FormatUint(uint64(stat.Gid), 10) == group.Gid
}

func prepareVelezKeyDir(_ context.Context, _ Options) (string, error) {
	err := velez.PrepareKeyDir(velez.DefaultKeyPath)
	if err != nil {
		return "", rerrors.Wrap(err, "error preparing velez key dir")
	}

	return velez.DefaultKeyPath + " ready for group " + vervGroup, nil
}

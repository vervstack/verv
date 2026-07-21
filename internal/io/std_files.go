package io

import (
	"os"

	"go.redsock.ru/rerrors"
)

const (
	DefaultDirPerm  os.FileMode = 0755
	DefaultFilePerm os.FileMode = 0644
)

func OverrideFile(pth string, content []byte) error {
	err := os.RemoveAll(pth)
	if err != nil {
		return rerrors.Wrap(err)
	}

	err = os.WriteFile(pth, content, DefaultFilePerm)
	if err != nil {
		return rerrors.Wrap(err)
	}

	return nil
}

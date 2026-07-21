package io

import (
	"os"

	"go.redsock.ru/rerrors"
)

func OverrideFile(pth string, content []byte) error {
	err := os.RemoveAll(pth)
	if err != nil {
		return rerrors.Wrap(err)
	}

	err = os.WriteFile(pth, content, 0755)
	if err != nil {
		return rerrors.Wrap(err)
	}

	return nil
}

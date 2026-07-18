package main_generators

import (
	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/internal/rw"
)

func GenerateMain() ([]byte, error) {
	out := &rw.RW{}

	err := mainTemplate.Execute(out, nil)
	if err != nil {
		return nil, rerrors.Wrap(err, "error generating main.go")
	}

	return out.Bytes(), nil
}

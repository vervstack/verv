package nodeclient

import "go.redsock.ru/rerrors"

var (
	ErrUnauthorized     = rerrors.New("node rejected credentials")
	ErrUnexpectedStatus = rerrors.New("node responded with unexpected status")
)

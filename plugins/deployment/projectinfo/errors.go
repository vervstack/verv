package projectinfo

import "go.redsock.ru/rerrors"

var (
	ErrServiceNameUnknown = rerrors.New("service name is unknown")
	ErrTagUnavailable     = rerrors.New("image tag is unavailable")

	errCommandFailed = rerrors.New("command failed")
)

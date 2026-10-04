package server

import "go.redsock.ru/rerrors"

var (
	errNotRoot          = rerrors.New("setup-server must run as root — use sudo")
	errUnsupportedOs    = rerrors.New("setup-server supports Ubuntu only")
	errEmptyUserName    = rerrors.New("user name must not be empty")
	errEmptyPassword    = rerrors.New("password must not be empty")
	errUnexpectedStatus = rerrors.New("unexpected http status")
	errUnknownCodename  = rerrors.New("os-release has no VERSION_CODENAME")

	errSysboxNoReleases     = rerrors.New("no numbered sysbox releases found")
	errSysboxVersionUnknown = rerrors.New("unknown sysbox version")
	errSysboxAssetMissing   = rerrors.New("sysbox release has no .deb package for this architecture")
	errSysboxRuntimeMissing = rerrors.New("sysbox-runc runtime is not registered in docker after install")
)

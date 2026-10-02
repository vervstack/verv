package upgrade

import (
	"go.redsock.ru/rerrors"
)

var (
	errUnsupportedTarget = rerrors.New("no release asset is published for this OS/architecture")
	errDownloadStatus    = rerrors.New("unexpected status downloading release asset")
)

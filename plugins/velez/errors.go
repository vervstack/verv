package velez

import "go.redsock.ru/rerrors"

var (
	errNoPortBinding = rerrors.New("existing velez container has no host port binding")
	errEmptyInspect  = rerrors.New("docker inspect returned no containers")
	errNoKeyMount    = rerrors.New("existing velez container has no keys mount")

	errKeyDirNeedsRoot = rerrors.New("cannot create /opt/velez without root — run as root or pass --key-path")
)

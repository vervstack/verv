package imagebuild

import "go.redsock.ru/rerrors"

var (
	errCommandFailed = rerrors.New("command failed")
	errBuildFailed   = rerrors.New("image build failed")
)

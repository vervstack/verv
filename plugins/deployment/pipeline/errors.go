package pipeline

import "go.redsock.ru/rerrors"

var (
	ErrNodePlatformUnknown = rerrors.New(
		"node did not report its os and arch - upgrade Velez on the node to a version that does")
	ErrNoRegistry          = rerrors.New("node has no registry configured")
	ErrEnvironmentRequired = rerrors.New("environment is required")
)

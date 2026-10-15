package nodes

import "go.redsock.ru/rerrors"

var (
	errEmptyName               = rerrors.New("node name is empty")
	errEmptyUrl                = rerrors.New("node url is empty")
	errEmptySecret             = rerrors.New("node secret is empty")
	errProjectScopeUnavailable = rerrors.New("project scope is unavailable: .verv directory does not exist")
	errUnknownScope            = rerrors.New("unknown node scope")
)

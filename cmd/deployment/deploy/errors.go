package deploy

import "go.redsock.ru/rerrors"

var (
	ErrNodeCredentialsMissing = rerrors.New(
		"CI mode requires environment variables VERV_NODE_URL and VERV_NODE_SECRET")

	errNoEnvironments = rerrors.New(
		"no environments found: create .verv/<env>/ or config/<env>.yaml to define one")
	errNodeUrlInvalid = rerrors.New("node address must start with http:// or https://")
	errNodeUrlEmpty   = rerrors.New("node address is empty")
	errSecretEmpty    = rerrors.New("secret is empty")
)

type rejectedNodeError string

func (e rejectedNodeError) Error() string {
	return string(e)
}

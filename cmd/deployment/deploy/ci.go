package deploy

import (
	"context"
	"os"
	"strings"

	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/plugins/deployment/nodeclient"
	"go.vervstack.ru/verv/plugins/deployment/pipeline"
)

const (
	envCi         = "CI"
	envNodeUrl    = "VERV_NODE_URL"
	envNodeSecret = "VERV_NODE_SECRET"
)

func isCiValue(value string) bool {
	normalized := strings.ToLower(strings.TrimSpace(value))

	return normalized != "" && normalized != "false" && normalized != "0"
}

func isCiMode() bool {
	return isCiValue(os.Getenv(envCi)) || !isStdinTerminal()
}

func isStdinTerminal() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	return info.Mode()&os.ModeCharDevice != 0
}

func (p *projectDeploy) runCi(ctx context.Context, params pipeline.Params) error {
	url := os.Getenv(envNodeUrl)
	secret := os.Getenv(envNodeSecret)

	if url == "" || secret == "" {
		return rerrors.Wrap(ErrNodeCredentialsMissing)
	}

	deps := p.newDeps(nodeclient.New(url, secret))

	err := pipeline.Run(ctx, deps, params)
	if err != nil {
		return rerrors.Wrap(err, "error running pipeline")
	}

	return nil
}

package deploy

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/internal/io"
	"go.vervstack.ru/verv/internal/io/colors"
	"go.vervstack.ru/verv/internal/menu"
	"go.vervstack.ru/verv/plugins/deployment/nodeclient"
	"go.vervstack.ru/verv/plugins/deployment/nodes"
)

const addNodeCheckTimeout = 5 * time.Second

func validateNodeUrl(value string) error {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return rerrors.Wrap(errNodeUrlEmpty)
	}

	if !strings.HasPrefix(trimmed, "http://") && !strings.HasPrefix(trimmed, "https://") {
		return rerrors.Wrap(errNodeUrlInvalid)
	}

	return nil
}

func validateSecret(value string) error {
	if strings.TrimSpace(value) == "" {
		return rerrors.Wrap(errSecretEmpty)
	}

	return nil
}

func nodeName(rawUrl string) string {
	parsed, err := url.Parse(rawUrl)
	if err != nil || parsed.Host == "" {
		return rawUrl
	}

	return parsed.Host
}

func (p *projectDeploy) addNode(ctx context.Context, store *nodes.Store) (aborted bool, err error) {
	var address, secret string

	addressInput := huh.NewInput().
		Title("Node address").
		Placeholder("https://velez.example.com").
		Validate(validateNodeUrl).
		Value(&address)

	secretInput := huh.NewInput().
		Title("Secret").
		EchoMode(huh.EchoModePassword).
		Validate(validateSecret).
		Value(&secret)

	form := huh.NewForm(huh.NewGroup(addressInput, secretInput)).WithShowHelp(false).WithTheme(menu.Theme())

	err = form.Run()
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return true, nil
		}

		return false, rerrors.Wrap(err, "error running add node form")
	}

	address = strings.TrimSpace(address)

	isReachable := p.checkNode(ctx, address, secret)
	if !isReachable {
		return false, nil
	}

	scope, aborted, err := p.chooseScope(store)
	if err != nil {
		return false, rerrors.Wrap(err, "error choosing node scope")
	}

	if aborted {
		return true, nil
	}

	node := nodes.Node{
		Name:   nodeName(address),
		Url:    address,
		Secret: secret,
	}

	err = store.Add(node, scope)
	if err != nil {
		return false, rerrors.Wrap(err, "error saving node")
	}

	p.io.PrintlnColored(colors.ColorGreen, fmt.Sprintf("%s Node added: %s", stepDone, node.Name))

	return false, nil
}

func (p *projectDeploy) checkNode(ctx context.Context, address, secret string) bool {
	checkCtx, cancel := context.WithTimeout(ctx, addNodeCheckTimeout)
	defer cancel()

	spinner := io.NewSpinner(p.io)
	spinner.Start("Checking node")

	_, err := nodeclient.New(address, secret).Version(checkCtx)
	if err != nil {
		spinner.Stop(false, "Checking node — failed")
		p.io.Error(colors.TerminalColor(colors.ColorRed) + fmt.Sprintf("%+v\n", err))

		return false
	}

	spinner.Stop(true, "Node is reachable")

	return true
}

func (p *projectDeploy) chooseScope(store *nodes.Store) (scope nodes.Scope, aborted bool, err error) {
	if !store.IsProjectScopeAvailable() {
		return nodes.ScopeGlobal, false, nil
	}

	sel := huh.NewSelect[nodes.Scope]().
		Title("Save for").
		Options(
			huh.NewOption("This project", nodes.ScopeProject),
			huh.NewOption("All projects", nodes.ScopeGlobal),
		).
		Value(&scope)

	form := huh.NewForm(huh.NewGroup(sel)).WithShowHelp(false).WithTheme(menu.Theme())

	err = form.Run()
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return nodes.ScopeGlobal, true, nil
		}

		return nodes.ScopeGlobal, false, rerrors.Wrap(err, "error running scope picker")
	}

	return scope, false, nil
}

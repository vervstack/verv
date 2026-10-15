package deploy

import (
	"errors"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"go.vervstack.ru/verv/plugins/deployment/nodeclient"
	"go.vervstack.ru/verv/plugins/deployment/nodes"
	"go.vervstack.ru/verv/plugins/deployment/pipeline"
)

const (
	scopeGlobalLabel  = "global"
	scopeProjectLabel = "project"

	reasonUnauthorized = "unauthorized"
	reasonUnreachable  = "unreachable"
	reasonNoRegistry   = "no registry"

	rejectUnauthorized = "node rejected the secret, fix it in the nodes file or add the node again"
	rejectUnreachable  = "node is unreachable"
	rejectNoRegistry   = "node has no registry, enable the registry plugin in Velez"
)

var disabledStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

func describeNode(entry nodes.Entry, status pipeline.Status) (label string, rejectReason string) {
	base := nodeLabel(entry, status)

	switch {
	case errors.Is(status.Err, nodeclient.ErrUnauthorized):
		return disabledStyle.Render(base + " — " + reasonUnauthorized), rejectUnauthorized
	case status.Err != nil:
		return disabledStyle.Render(base + " — " + reasonUnreachable), rejectUnreachable
	case !status.IsRegistryEnabled:
		return disabledStyle.Render(base + " — " + reasonNoRegistry), rejectNoRegistry
	}

	return base, ""
}

func nodeLabel(entry nodes.Entry, status pipeline.Status) string {
	parts := []string{entry.Name}

	if status.Platform.Version != "" {
		parts = append(parts, status.Platform.Version)
	}

	parts = append(parts, "("+scopeLabel(entry.Scope)+")")

	return strings.Join(parts, "  ")
}

func scopeLabel(scope nodes.Scope) string {
	if scope == nodes.ScopeProject {
		return scopeProjectLabel
	}

	return scopeGlobalLabel
}

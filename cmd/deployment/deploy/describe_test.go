package deploy

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/plugins/deployment/nodeclient"
	"go.vervstack.ru/verv/plugins/deployment/nodes"
	"go.vervstack.ru/verv/plugins/deployment/pipeline"
)

func Test_DescribeNode_Statuses(t *testing.T) {
	entry := nodes.Entry{
		Node:  nodes.Node{Name: "prod-1", Url: "https://prod-1", Secret: "s"},
		Scope: nodes.ScopeProject,
	}

	cases := []struct {
		name       string
		status     pipeline.Status
		wantLabel  string
		wantSuffix string
		wantReject string
	}{
		{
			name: "ok",
			status: pipeline.Status{
				Platform:          nodeclient.Platform{Version: "v1.2.3"},
				IsRegistryEnabled: true,
			},
			wantLabel: "prod-1  v1.2.3  (project)",
		},
		{
			name: "unauthorized",
			status: pipeline.Status{
				Err: rerrors.Wrap(nodeclient.ErrUnauthorized, "status: 403"),
			},
			wantLabel:  "prod-1  (project)",
			wantSuffix: "— unauthorized",
			wantReject: rejectUnauthorized,
		},
		{
			name:       "unreachable",
			status:     pipeline.Status{Err: rerrors.New("connection refused")},
			wantLabel:  "prod-1  (project)",
			wantSuffix: "— unreachable",
			wantReject: rejectUnreachable,
		},
		{
			name: "no registry",
			status: pipeline.Status{
				Platform: nodeclient.Platform{Version: "v1.2.3"},
			},
			wantLabel:  "prod-1  v1.2.3  (project)",
			wantSuffix: "— no registry",
			wantReject: rejectNoRegistry,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			label, reject := describeNode(entry, tc.status)

			require.Contains(t, label, tc.wantLabel)
			require.Contains(t, label, tc.wantSuffix)
			require.Equal(t, tc.wantReject, reject)
		})
	}
}

func Test_NodeName_HostOrRaw(t *testing.T) {
	require.Equal(t, "velez.example.com:8080", nodeName("https://velez.example.com:8080/api"))
	require.Equal(t, "not a url", nodeName("not a url"))
}

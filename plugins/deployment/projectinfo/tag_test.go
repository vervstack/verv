package projectinfo_test

import (
	"context"
	"strings"
	"testing"

	"go.redsock.ru/rerrors"

	"github.com/stretchr/testify/require"

	"go.vervstack.ru/verv/plugins/deployment/projectinfo"
)

var errFake = rerrors.New("fake failure")

type fakeResponse struct {
	stdout string
	err    error
}

type fakeExecutor struct {
	responses map[string]fakeResponse
}

func (f *fakeExecutor) Run(_ context.Context, _ string, tool string, args ...string) (string, error) {
	key := tool + " " + strings.Join(args, " ")

	response, ok := f.responses[key]
	if !ok {
		return "", errFake
	}

	return response.stdout, response.err
}

func Test_Tag_Variants(t *testing.T) {
	cases := []struct {
		name      string
		responses map[string]fakeResponse
		want      string
		isFailure bool
	}{
		{
			name: "clean tree",
			responses: map[string]fakeResponse{
				"git rev-parse --short HEAD": {stdout: "abc1234\n"},
				"git status --porcelain":     {stdout: ""},
			},
			want: "abc1234",
		},
		{
			name: "dirty tree",
			responses: map[string]fakeResponse{
				"git rev-parse --short HEAD": {stdout: "abc1234\n"},
				"git status --porcelain":     {stdout: " M main.go\n"},
			},
			want: "abc1234-dirty",
		},
		{
			name: "not a git repository",
			responses: map[string]fakeResponse{
				"git rev-parse --short HEAD": {err: errFake},
			},
			isFailure: true,
		},
		{
			name: "status fails",
			responses: map[string]fakeResponse{
				"git rev-parse --short HEAD": {stdout: "abc1234\n"},
				"git status --porcelain":     {err: errFake},
			},
			isFailure: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			executor := &fakeExecutor{responses: tc.responses}

			got, err := projectinfo.Tag(context.Background(), executor, "/project")

			if tc.isFailure {
				require.ErrorIs(t, err, projectinfo.ErrTagUnavailable)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}

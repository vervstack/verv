package deploy

import (
	"testing"

	"github.com/stretchr/testify/require"

	"go.vervstack.ru/verv/internal/io"
	"go.vervstack.ru/verv/plugins/deployment/pipeline"
)

func Test_IsCiValue_Values(t *testing.T) {
	cases := []struct {
		value string
		want  bool
	}{
		{"", false},
		{"false", false},
		{"FALSE", false},
		{"0", false},
		{" 0 ", false},
		{"true", true},
		{"1", true},
		{"yes", true},
	}

	for _, tc := range cases {
		t.Run(tc.value, func(t *testing.T) {
			require.Equal(t, tc.want, isCiValue(tc.value))
		})
	}
}

func Test_RunCi_CredentialsMissing(t *testing.T) {
	t.Setenv(envNodeUrl, "")
	t.Setenv(envNodeSecret, "secret")

	proc := projectDeploy{io: io.StdIO{}, wd: t.TempDir()}

	err := proc.runCi(t.Context(), pipeline.Params{Environment: "prod"})

	require.ErrorIs(t, err, ErrNodeCredentialsMissing)
}

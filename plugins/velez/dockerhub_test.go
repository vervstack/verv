package velez

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ParseTags_Scenarios(t *testing.T) {
	cases := []struct {
		name    string
		body    string
		want    []string
		wantErr bool
	}{
		{
			name: "three results",
			body: `{"results":[{"name":"v0.1.90"},{"name":"v0.1.89"},{"name":"v0.1.88"}]}`,
			want: []string{"v0.1.90", "v0.1.89", "v0.1.88"},
		},
		{
			name:    "empty results",
			body:    `{"results":[]}`,
			wantErr: true,
		},
		{
			name: "latest entry is skipped",
			body: `{"results":[{"name":"latest"},{"name":"v0.1.90"},{"name":"v0.1.89"}]}`,
			want: []string{"v0.1.90", "v0.1.89"},
		},
		{
			name: "results beyond max are truncated",
			body: `{"results":[{"name":"v4"},{"name":"v3"},{"name":"v2"},{"name":"v1"}]}`,
			want: []string{"v4", "v3", "v2"},
		},
		{
			name:    "invalid json",
			body:    `not json`,
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseTags([]byte(tc.body))
			if tc.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}

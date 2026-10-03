package server

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_MergeAuthorizedKeys_Scenarios(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		existing string
		fetched  string
		want     string
	}{
		{"empty existing", "", "key-a\nkey-b\n", "key-a\nkey-b\n"},
		{"skips present keys", "key-a\n", "key-a\nkey-b\n", "key-a\nkey-b\n"},
		{"adds newline to existing", "key-a", "key-b", "key-a\nkey-b\n"},
		{"ignores blank lines and whitespace", "key-a\n", "\n  key-a  \n\n key-b\n", "key-a\nkey-b\n"},
		{"deduplicates fetched", "", "key-a\nkey-a\n", "key-a\n"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := mergeAuthorizedKeys(tc.existing, tc.fetched)

			require.Equal(t, tc.want, got)
		})
	}
}

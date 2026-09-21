package folder

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_isUnchangedFromOlderVersion(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		older   []byte
		current []byte
		want    bool
	}{
		{"identical single byte", []byte("a"), []byte("a"), true},
		{"changed single byte", []byte("a"), []byte("b"), false},
		{"identical multi byte", []byte("abcdef"), []byte("abcdef"), true},
		{"changed first byte", []byte("abcdef"), []byte("xbcdef"), false},
		{"changed last byte", []byte("abcdef"), []byte("abcdex"), false},
		{"changed middle byte", []byte("abcdef"), []byte("abXdef"), false},
		{"both empty", []byte{}, []byte{}, true},
		{"different length", []byte("ab"), []byte("abc"), false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			f := &Folder{
				olderVersion: tc.older,
				Content:      tc.current,
			}

			require.Equal(t, tc.want, f.isUnchangedFromOlderVersion())
		})
	}
}

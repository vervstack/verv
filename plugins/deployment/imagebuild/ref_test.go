package imagebuild_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"go.vervstack.ru/verv/plugins/deployment/imagebuild"
)

func Test_ImageRef_Variants(t *testing.T) {
	cases := []struct {
		name        string
		registryUrl string
		want        string
	}{
		{"bare host", "registry.example.com", "registry.example.com/artel:abc"},
		{"https scheme", "https://registry.example.com", "registry.example.com/artel:abc"},
		{"http scheme with port", "http://registry.example.com:5000", "registry.example.com:5000/artel:abc"},
		{"trailing slash", "https://registry.example.com/", "registry.example.com/artel:abc"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := imagebuild.ImageRef(tc.registryUrl, "artel", "abc")

			require.Equal(t, tc.want, got)
		})
	}
}

package middleware_generators

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_GenerateCookieNames(t *testing.T) {
	t.Parallel()

	got := GenerateCookieNames()
	require.Equal(t, cookieNamesFile, got)
}

func Test_GenerateCookieAnnotator(t *testing.T) {
	t.Parallel()

	got := GenerateCookieAnnotator()
	require.Equal(t, cookieAnnotatorFile, got)
}

func Test_GenerateCookieResponse(t *testing.T) {
	t.Parallel()

	got := GenerateCookieResponse()
	require.Equal(t, cookieResponseFile, got)
}

func Test_GenerateCSRFInterceptor(t *testing.T) {
	t.Parallel()

	got := GenerateCSRFInterceptor()
	require.Equal(t, csrfInterceptorFile, got)
}

func Test_GenerateRequestSchemeAnnotator(t *testing.T) {
	t.Parallel()

	got := GenerateRequestSchemeAnnotator()
	require.Equal(t, requestSchemeAnnotatorFile, got)
}

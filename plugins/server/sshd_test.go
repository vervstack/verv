package server

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_EnablePubkeyAuth_UncommentsDirective(t *testing.T) {
	t.Parallel()

	got := enablePubkeyAuth("Port 22\n#PubkeyAuthentication yes\nUsePAM yes\n")

	require.Equal(t, "Port 22\nPubkeyAuthentication yes\nUsePAM yes\n", got)
}

func Test_EnablePubkeyAuth_IsIdempotent(t *testing.T) {
	t.Parallel()

	once := enablePubkeyAuth("#PubkeyAuthentication yes\n")

	require.Equal(t, once, enablePubkeyAuth(once))
}

func Test_EnablePubkeyAuth_LeavesOtherContentUntouched(t *testing.T) {
	t.Parallel()

	content := "#PubkeyAuthentication no\nPasswordAuthentication yes\n"

	require.Equal(t, content, enablePubkeyAuth(content))
}

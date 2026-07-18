package main_generators

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_GenerateMain(t *testing.T) {
	got, err := GenerateMain()
	require.NoError(t, err)
	require.Equal(t, mainPattern, string(got))
}

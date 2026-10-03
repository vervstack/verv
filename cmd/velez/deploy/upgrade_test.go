package deploy

import (
	"testing"

	"github.com/stretchr/testify/require"

	"go.vervstack.ru/verv/internal/processor"
)

func Test_NewUpgradeCommand_StartsHiddenWithoutPortFlag(t *testing.T) {
	c := NewUpgradeCommand(processor.New())

	require.True(t, c.Hidden)
	require.Nil(t, c.Flags().Lookup(PortFlag))
	require.Nil(t, c.Flags().Lookup(KeyPathFlag))
	require.NotNil(t, c.Flags().Lookup(DebugFlag))
}

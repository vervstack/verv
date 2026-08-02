package menu_test

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"go.vervstack.ru/verv/internal/menu"
)

func noopRun(*cobra.Command, []string) {}

func Test_BuildEntries_FiltersHiddenAndDeprecated(t *testing.T) {
	t.Parallel()

	initCmd := &cobra.Command{
		Use:         "init",
		Short:       "Initialize a new project",
		Long:        "Initialize a new project with an interactive name prompt",
		Run:         noopRun,
		Annotations: map[string]string{"verv:emoji": "🚀"},
	}
	hiddenCmd := &cobra.Command{
		Use:    "hidden",
		Short:  "should not appear",
		Hidden: true,
		Run:    noopRun,
	}
	deprecatedCmd := &cobra.Command{
		Use:        "old",
		Short:      "should not appear",
		Deprecated: "use something else",
		Run:        noopRun,
	}
	tidyCmd := &cobra.Command{
		Use:         "tidy",
		Short:       "Tidy up a project",
		Long:        "Run the tidy pipeline over the current project",
		Run:         noopRun,
		Annotations: map[string]string{"verv:emoji": "🧹", "verv:requiresProject": "true"},
	}

	entries := menu.BuildEntries([]*cobra.Command{initCmd, hiddenCmd, deprecatedCmd, tidyCmd})

	require.Len(t, entries, 2)

	require.Equal(t, "init", entries[0].Name)
	require.Equal(t, initCmd.Short, entries[0].Short)
	require.Equal(t, initCmd.Long, entries[0].Long)
	require.Same(t, initCmd, entries[0].Cmd)
	require.Equal(t, "🚀", entries[0].Emoji)
	require.False(t, entries[0].RequiresProject)

	require.Equal(t, "tidy", entries[1].Name)
	require.Equal(t, tidyCmd.Short, entries[1].Short)
	require.Equal(t, tidyCmd.Long, entries[1].Long)
	require.Same(t, tidyCmd, entries[1].Cmd)
	require.Equal(t, "🧹", entries[1].Emoji)
	require.True(t, entries[1].RequiresProject)
}

func Test_BuildEntries_DefaultEmojiWhenUnset(t *testing.T) {
	t.Parallel()

	cmd := &cobra.Command{Use: "noannotations", Run: noopRun}

	entries := menu.BuildEntries([]*cobra.Command{cmd})

	require.Len(t, entries, 1)
	require.NotEmpty(t, entries[0].Emoji)
	require.False(t, entries[0].RequiresProject)
}

func Test_BuildEntries_EmptyInput(t *testing.T) {
	t.Parallel()

	entries := menu.BuildEntries(nil)

	require.NotNil(t, entries)
	require.Empty(t, entries)
}

func Test_BuildEntries_PreservesOrder(t *testing.T) {
	t.Parallel()

	cmds := []*cobra.Command{
		{Use: "add", Run: noopRun},
		{Use: "init", Run: noopRun},
		{Use: "tidy", Run: noopRun},
	}

	entries := menu.BuildEntries(cmds)

	require.Len(t, entries, 3)
	require.Equal(t, "add", entries[0].Name)
	require.Equal(t, "init", entries[1].Name)
	require.Equal(t, "tidy", entries[2].Name)
}

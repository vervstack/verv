package menu

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func Test_OrderedGroups_Scenarios(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name          string
		headlessLinux bool
		want          []string
	}{
		{"headless linux leads with velez", true, []string{GroupVelez, GroupProject}},
		{"otherwise leads with project", false, []string{GroupProject, GroupVelez}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := orderedGroups(tc.headlessLinux)

			require.Equal(t, tc.want, got)
		})
	}
}

func Test_GroupedOptions_InsertsHeaderPerNonEmptySection(t *testing.T) {
	t.Parallel()

	tidyCmd := &cobra.Command{Use: "tidy"}
	deployCmd := &cobra.Command{Use: "deploy-velez"}

	entries := []Entry{
		{Name: "tidy", Cmd: tidyCmd, Group: GroupProject},
		{Name: "deploy-velez", Cmd: deployCmd, Group: GroupVelez},
	}

	options := groupedOptions(entries, []string{GroupProject, GroupVelez}, Header{})

	require.Len(t, options, 4)
	require.Nil(t, options[0].Value)
	require.Same(t, tidyCmd, options[1].Value)
	require.Nil(t, options[2].Value)
	require.Same(t, deployCmd, options[3].Value)
}

func Test_GroupedOptions_SkipsEmptySection(t *testing.T) {
	t.Parallel()

	tidyCmd := &cobra.Command{Use: "tidy"}
	entries := []Entry{{Name: "tidy", Cmd: tidyCmd, Group: GroupProject}}

	options := groupedOptions(entries, []string{GroupVelez, GroupProject}, Header{})

	require.Len(t, options, 2)
	require.Nil(t, options[0].Value)
	require.Same(t, tidyCmd, options[1].Value)
}

func Test_FirstSelectable_SkipsHeaders(t *testing.T) {
	t.Parallel()

	tidyCmd := &cobra.Command{Use: "tidy"}
	options := []huh.Option[*cobra.Command]{
		huh.NewOption("Project", (*cobra.Command)(nil)),
		huh.NewOption("tidy", tidyCmd),
	}

	got := firstSelectable(options)

	require.Same(t, tidyCmd, got)
}

func Test_GroupedOptions_HeaderPlainAndOptionsIndented(t *testing.T) {
	t.Parallel()

	tidyCmd := &cobra.Command{Use: "tidy", Short: "🧹 tidy"}
	entries := []Entry{{Name: "tidy", Emoji: "🧹", Cmd: tidyCmd, Group: GroupProject}}

	options := groupedOptions(entries, []string{GroupProject}, Header{})

	require.Len(t, options, 2)
	require.NotContains(t, options[0].Key, "─")
	require.True(t, strings.HasPrefix(options[1].Key, optionIndent))
}

func Test_IsSelectNavKey_Scenarios(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		msg  tea.Msg
		want bool
	}{
		{"up arrow", tea.KeyMsg{Type: tea.KeyUp}, true},
		{"down arrow", tea.KeyMsg{Type: tea.KeyDown}, true},
		{"vim j", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")}, true},
		{"vim k", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")}, true},
		{"enter is not nav", tea.KeyMsg{Type: tea.KeyEnter}, false},
		{"non-key message", tea.WindowSizeMsg{}, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := isSelectNavKey(tc.msg)

			require.Equal(t, tc.want, got)
		})
	}
}

func Test_NavModel_SkipsHeaderOnDown(t *testing.T) {
	t.Parallel()

	optA := &cobra.Command{Use: "optA"}
	optB := &cobra.Command{Use: "optB"}

	// Mirrors production Select(): chosen is preset to a real command before
	// Value() binds it, so the field's initial cursor lands on optA — never
	// on the header — exactly like firstSelectable does for the real picker.
	chosen := optA

	sel := huh.NewSelect[*cobra.Command]().
		Options(
			huh.NewOption("optA", optA),
			huh.NewOption("Section", (*cobra.Command)(nil)),
			huh.NewOption("optB", optB),
		).
		Value(&chosen)

	form := huh.NewForm(huh.NewGroup(sel))
	form.SubmitCmd = tea.Quit
	form.CancelCmd = tea.Interrupt

	require.Same(t, optA, chosen)

	model := navModel{
		form:           form,
		sel:            sel,
		chosen:         &chosen,
		descriptionFor: func() string { return "" },
		maxSkip:        3,
	}

	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})

	navUpdated, ok := updated.(navModel)
	require.True(t, ok)
	require.Same(t, optB, *navUpdated.chosen)
}

func Test_NavModel_UpdatesDescriptionOnEveryNavStep(t *testing.T) {
	t.Parallel()

	optA := &cobra.Command{Use: "optA"}
	optB := &cobra.Command{Use: "optB"}

	descByCmd := map[*cobra.Command]string{
		optA: "description A",
		optB: "description B",
	}

	chosen := optA

	sel := huh.NewSelect[*cobra.Command]().
		Options(huh.NewOption("optA", optA), huh.NewOption("optB", optB)).
		Value(&chosen)

	form := huh.NewForm(huh.NewGroup(sel))
	form.SubmitCmd = tea.Quit
	form.CancelCmd = tea.Interrupt

	model := navModel{
		form:           form,
		sel:            sel,
		chosen:         &chosen,
		descriptionFor: func() string { return descByCmd[chosen] },
		maxSkip:        2,
	}

	// Asserting on model.View() (what tea.Program actually renders), not
	// sel.View() directly — Group.View() replays a viewport cached inside
	// form.Update, so a description pushed after that call needs a rebuild
	// to actually show up here. Asserting on sel.View() would pass even
	// without that rebuild and miss the bug entirely.
	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})

	navUpdated := updated.(navModel)
	require.Contains(t, navUpdated.View(), "description B")
	require.NotContains(t, navUpdated.View(), "description A")
}

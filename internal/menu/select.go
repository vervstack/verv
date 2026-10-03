package menu

import (
	"errors"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"go.redsock.ru/rerrors"
)

const (
	disabledSuffix       = " (needs a verv project)"
	dockerDisabledSuffix = " (needs docker engine)"
	optionIndent         = "  "

	wordmark = `
██╗   ██╗███████╗██████╗ ██╗   ██╗
██║   ██║██╔════╝██╔══██╗██║   ██║
██║   ██║█████╗  ██████╔╝██║   ██║
╚██╗ ██╔╝██╔══╝  ██╔══██╗╚██╗ ██╔╝
 ╚████╔╝ ███████╗██║  ██║ ╚████╔╝
  ╚═══╝  ╚══════╝╚═╝  ╚═╝  ╚═══╝  `
)

var (
	errGroupHeaderNotSelectable = rerrors.New("that's a section header — pick an option below it")
	errDockerRequired           = rerrors.New(
		"Can't run without docker engine - install via setup server option or install docker manually")

	// groupLabels names each picker section header, in GroupProject's default
	// display order (headless Linux swaps this — see orderedGroups).
	groupLabels = map[string]string{
		GroupProject:  "Project",
		GroupVelez:    "Velez node",
		GroupSettings: "Settings",
	}

	headerStyle   = lipgloss.NewStyle().Bold(true)
	selectorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).SetString("◆─ ")
	wordmarkStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212"))
	badgeStyle    = lipgloss.NewStyle().Bold(true)
	pathStyle     = lipgloss.NewStyle().Faint(true)
	updateStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("214"))
)

// Theme returns huh's default (Charm) theme with the selected-row cursor
// swapped for a diamond-and-dash marker instead of the plain "> " arrow.
func Theme() *huh.Theme {
	t := huh.ThemeCharm()

	t.Focused.SelectSelector = selectorStyle
	t.Blurred.SelectSelector = selectorStyle

	return t
}

func updateNotice(tag string) string {
	return "⬆ Update available: " + tag + " — pick Settings > upgrade"
}

// Select renders an arrow-key-navigable picker over entries and returns the
// chosen Entry. It returns (nil, nil) if the user aborts (Ctrl+C/Esc) — a
// quiet abort rather than an error — and (nil, err) on any other failure.
func Select(entries []Entry, header Header) (*Entry, error) {
	options := groupedOptions(entries, orderedGroups(runningHeadlessLinux()), header)

	chosen := firstSelectable(options)

	descriptionFor := func() string {
		e := findEntry(entries, chosen)
		if e == nil {
			return ""
		}

		return e.Long
	}

	sel := huh.NewSelect[*cobra.Command]().
		Title("What would you like to do?").
		Options(options...).
		Value(&chosen).
		Description(descriptionFor()).
		Validate(validateChoice(entries, header))

	form := huh.NewForm(huh.NewGroup(sel)).
		WithShowHelp(false).
		WithKeyMap(pickerKeyMap()).
		WithTheme(Theme())

	// Run our own tea.Program instead of form.Run(), wrapped in navModel, so
	// that Up/Down (and their j/k/ctrl+n/ctrl+p aliases) step over the
	// unselectable group-header options instead of landing on them —
	// something huh's Select has no built-in support for.
	form.SubmitCmd = tea.Quit
	form.CancelCmd = tea.Interrupt

	model := navModel{
		form:           form,
		sel:            sel,
		chosen:         &chosen,
		descriptionFor: descriptionFor,
		maxSkip:        len(options),
		header:         renderHeader(header),
	}

	_, err := tea.NewProgram(model, tea.WithOutput(os.Stderr)).Run()
	if err != nil {
		if errors.Is(err, tea.ErrInterrupted) {
			return nil, nil
		}

		return nil, rerrors.Wrap(err, "error running command picker")
	}

	if form.State == huh.StateAborted {
		return nil, nil
	}

	return findEntry(entries, chosen), nil
}

func findEntry(entries []Entry, cmd *cobra.Command) *Entry {
	for i := range entries {
		if entries[i].Cmd == cmd {
			return &entries[i]
		}
	}

	return nil
}

func validateChoice(entries []Entry, header Header) func(*cobra.Command) error {
	return func(cmd *cobra.Command) error {
		if cmd == nil {
			return rerrors.Wrap(errGroupHeaderNotSelectable)
		}

		e := findEntry(entries, cmd)
		if e != nil && e.RequiresProject && !header.IsVervProject {
			return rerrors.New(e.Name + " needs an existing verv project — run `verv init` first")
		}

		if e != nil && e.RequiresDocker && !header.IsDockerInstalled {
			return rerrors.Wrap(errDockerRequired)
		}

		return nil
	}
}

func pickerKeyMap() *huh.KeyMap {
	km := huh.NewDefaultKeyMap()
	km.Quit.SetKeys("ctrl+c", "esc")
	km.Select.Filter.SetEnabled(false)
	km.Select.GotoTop.SetEnabled(false)
	km.Select.GotoBottom.SetEnabled(false)
	km.Select.HalfPageUp.SetEnabled(false)
	km.Select.HalfPageDown.SetEnabled(false)

	return km
}

// asForm narrows a tea.Model returned by huh.Form.Update back to the form,
// keeping the previous form if huh ever returns something else.
func asForm(model tea.Model, fallback *huh.Form) *huh.Form {
	f, ok := model.(*huh.Form)
	if !ok {
		return fallback
	}

	return f
}

// navModel wraps a *huh.Form to work around two gaps in huh's Select field:
//
//   - Cursor-move keys have no way to skip a non-selectable item, so a
//     group-header option (nil command) would otherwise be landable. This
//     repeats the same key internally until the cursor reaches a real entry.
//   - DescriptionFunc re-evaluates only when a hash of its bindings changes,
//     and that hash is computed by reflecting into the bound value — which
//     fails (and is silently swallowed) for anything containing a func field,
//     as every *cobra.Command does. The description would then never update
//     past its first render. Description is instead pushed by hand here,
//     after every keystroke, bypassing that mechanism entirely.
//
// It also renders the header itself, rather than the caller Println-ing it
// before the program starts: bubbletea's inline (non-alt-screen) renderer
// repaints its own view region every frame, so a header printed beforehand
// would get overwritten by the first repaint unless it's part of the same
// View().
type navModel struct {
	form           *huh.Form
	sel            *huh.Select[*cobra.Command]
	chosen         **cobra.Command
	descriptionFor func() string
	maxSkip        int
	header         string
}

func (m navModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m navModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !isSelectNavKey(msg) {
		newForm, cmd := m.form.Update(msg)

		m.form = asForm(newForm, m.form)

		return m, cmd
	}

	var cmd tea.Cmd

	for range m.maxSkip {
		var newForm tea.Model

		newForm, cmd = m.form.Update(msg)
		m.form = asForm(newForm, m.form)

		if *m.chosen != nil || m.form.State != huh.StateNormal {
			break
		}
	}

	m.sel.Description(m.descriptionFor())

	// Group caches its rendered viewport inside form.Update and View just
	// replays that cache — so the push above lands after this frame's cache
	// was already built and would only show up one keystroke late. Drive one
	// more (unhandled, no-op) message through to force a rebuild against the
	// description we just set.
	newForm, _ := m.form.Update(refreshMsg{})

	m.form = asForm(newForm, m.form)

	return m, cmd
}

// refreshMsg is never matched by huh's Select/Group Update switches, so
// routing it through the form is a no-op except for the unconditional
// Group.buildView() call at the end of Group.Update — used to force a
// viewport rebuild after mutating field state directly (see above).
type refreshMsg struct{}

func (m navModel) View() string {
	return m.header + "\n\n" + m.form.View()
}

// isSelectNavKey reports whether msg is one of the Select field's cursor-move
// keys — the only ones that can land on a group-header option.
func isSelectNavKey(msg tea.Msg) bool {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return false
	}

	switch keyMsg.String() {
	case "up", "down", "k", "j", "ctrl+k", "ctrl+j", "ctrl+p", "ctrl+n":
		return true
	default:
		return false
	}
}

// orderedGroups returns the picker section order. On a headless Linux
// machine — the profile of a server someone is deploying a Velez node onto —
// the Velez section leads; everywhere else Project leads. Settings always comes last.
func orderedGroups(headlessLinux bool) []string {
	if headlessLinux {
		return []string{GroupVelez, GroupProject, GroupSettings}
	}

	return []string{GroupProject, GroupVelez, GroupSettings}
}

// groupedOptions renders entries as huh options, section by section in the
// given group order. Each non-empty section is preceded by a header option
// bound to a nil command — Select's Validate rejects choosing it, so it acts
// as an unselectable divider without needing huh to support disabled items.
func groupedOptions(entries []Entry, groups []string, header Header) []huh.Option[*cobra.Command] {
	options := make([]huh.Option[*cobra.Command], 0, len(entries)+len(groups))

	for _, group := range groups {
		var groupOptions []huh.Option[*cobra.Command]

		for _, e := range entries {
			if e.Group != group {
				continue
			}

			label := optionIndent + e.Emoji + " " + e.Name
			if e.RequiresProject && !header.IsVervProject {
				label += disabledSuffix
			}

			if e.RequiresDocker && !header.IsDockerInstalled {
				label += dockerDisabledSuffix
			}

			groupOptions = append(groupOptions, huh.NewOption(label, e.Cmd))
		}

		if len(groupOptions) == 0 {
			continue
		}

		options = append(options, huh.NewOption(headerStyle.Render(groupLabels[group]), (*cobra.Command)(nil)))
		options = append(options, groupOptions...)
	}

	return options
}

// firstSelectable returns the first non-header option's value, so the picker
// opens with the cursor on a real command rather than a divider.
func firstSelectable(options []huh.Option[*cobra.Command]) *cobra.Command {
	for _, o := range options {
		if o.Value != nil {
			return o.Value
		}
	}

	return nil
}

// renderHeader renders the ASCII "VERV" wordmark, followed by a
// project-aware status line: `<emoji> <name> · v<version> · <path>` when
// standing inside a verv project (detection is strict — only the .verv
// marker counts), or just the plain working directory path otherwise.
func renderHeader(header Header) string {
	lines := []string{wordmarkStyle.Render(wordmark)}

	if header.UpdateVersion != "" {
		lines = append(lines, updateStyle.Render(updateNotice(header.UpdateVersion)))
	}

	if !header.IsVervProject {
		lines = append(lines, pathStyle.Render(header.Path))

		return strings.Join(lines, "\n")
	}

	line := header.Emoji + " " + badgeStyle.Render(header.Name)
	if header.Version != "" {
		line += " · " + versionLabel(header.Version)
	}

	line += " · " + header.Path

	lines = append(lines, line)

	return strings.Join(lines, "\n")
}

// versionLabel renders a version string with a "v" prefix, tolerating
// versions that already carry one (e.g. "v0.0.1" from AppInfo.Version) so we
// never print "vv0.0.1".
func versionLabel(version string) string {
	if strings.HasPrefix(version, "v") || strings.HasPrefix(version, "V") {
		return version
	}

	return "v" + version
}

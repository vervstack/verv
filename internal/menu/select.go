package menu

import (
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"go.redsock.ru/rerrors"
)

const disabledSuffix = " (needs a verv project)"

const wordmark = `
██╗   ██╗███████╗██████╗ ██╗   ██╗
██║   ██║██╔════╝██╔══██╗██║   ██║
██║   ██║█████╗  ██████╔╝██║   ██║
╚██╗ ██╔╝██╔══╝  ██╔══██╗╚██╗ ██╔╝
 ╚████╔╝ ███████╗██║  ██║ ╚████╔╝
  ╚═══╝  ╚══════╝╚═╝  ╚═╝  ╚═══╝  `

var (
	wordmarkStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212"))
	badgeStyle    = lipgloss.NewStyle().Bold(true)
	pathStyle     = lipgloss.NewStyle().Faint(true)
)

// Select renders an arrow-key-navigable picker over entries and returns the
// chosen Entry. It returns (nil, nil) if the user aborts (Ctrl+C/Esc) — a
// quiet abort rather than an error — and (nil, err) on any other failure.
func Select(entries []Entry, header Header) (*Entry, error) {
	printHeader(header)

	options := make([]huh.Option[*cobra.Command], 0, len(entries))
	for _, e := range entries {
		label := e.Emoji + " " + e.Name
		if e.RequiresProject && !header.IsVervProject {
			label += disabledSuffix
		}

		options = append(options, huh.NewOption(label, e.Cmd))
	}

	var chosen *cobra.Command

	entryFor := func(cmd *cobra.Command) *Entry {
		for i := range entries {
			if entries[i].Cmd == cmd {
				return &entries[i]
			}
		}

		return nil
	}

	descriptionFor := func() string {
		e := entryFor(chosen)
		if e == nil {
			return ""
		}

		return e.Long
	}

	sel := huh.NewSelect[*cobra.Command]().
		Title("What would you like to do?").
		Options(options...).
		Value(&chosen).
		DescriptionFunc(descriptionFor, &chosen).
		Validate(func(cmd *cobra.Command) error {
			e := entryFor(cmd)
			if e != nil && e.RequiresProject && !header.IsVervProject {
				return rerrors.New(e.Name + " needs an existing verv project — run `verv init` first")
			}

			return nil
		})

	km := huh.NewDefaultKeyMap()
	km.Quit.SetKeys("ctrl+c", "esc")

	form := huh.NewForm(huh.NewGroup(sel)).WithShowHelp(false).WithKeyMap(km)

	err := form.Run()
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return nil, nil
		}

		return nil, rerrors.Wrap(err, "error running command picker")
	}

	for i := range entries {
		if entries[i].Cmd == chosen {
			return &entries[i], nil
		}
	}

	return nil, nil
}

// printHeader renders the ASCII "VERV" wordmark, followed by a project-aware
// status line: `<emoji> <name> · v<version> · <path>` when standing inside a
// verv project (detection is strict — only the .verv marker counts), or just
// the plain working directory path otherwise.
func printHeader(header Header) {
	fmt.Println(wordmarkStyle.Render(wordmark))

	if !header.IsVervProject {
		fmt.Println(pathStyle.Render(header.Path))
		fmt.Println()

		return
	}

	line := header.Emoji + " " + badgeStyle.Render(header.Name)
	if header.Version != "" {
		line += " · " + versionLabel(header.Version)
	}

	line += " · " + header.Path

	fmt.Println(line)
	fmt.Println()
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

package menu

import (
	"github.com/spf13/cobra"
)

const (
	// annotationEmoji, annotationRequiresProject, and annotationGroup are the
	// cobra.Command Annotations keys commands set to opt into a picker icon,
	// a "needs an existing verv project" gate, and/or a picker section,
	// keeping BuildEntries free of a hardcoded command-name list.
	annotationEmoji           = "verv:emoji"
	annotationRequiresProject = "verv:requiresProject"
	annotationGroup           = "verv:group"
	annotationRequiresDocker  = "verv:requiresDocker"

	// annotationShowsSetupProgress marks the command whose picker label shows how many
	// setup steps are complete.
	annotationShowsSetupProgress = "verv:showsSetupProgress"

	annotationTrue = "true"

	defaultEntryEmoji = "▫️"

	// GroupProject, GroupVelez and GroupSettings are the picker section keys. A command with
	// no verv:group annotation defaults to GroupProject.
	GroupProject  = "project"
	GroupVelez    = "velez"
	GroupSettings = "settings"
)

// Entry describes a single top-level command as presented in the interactive
// command picker.
type Entry struct {
	Name            string
	Short           string
	Long            string
	Cmd             *cobra.Command
	Emoji           string
	RequiresProject bool
	RequiresDocker  bool
	Group           string

	ShowsSetupProgress bool
}

// BuildEntries converts cobra commands into menu Entry values, filtering out
// anything cobra itself would not surface as an "available" command (hidden,
// deprecated, or non-runnable) — the same predicate cobra's own help
// template relies on. This keeps the menu in sync with whatever commands are
// registered on root, with no hardcoded name list.
func BuildEntries(cmds []*cobra.Command) []Entry {
	entries := make([]Entry, 0, len(cmds))

	for _, c := range cmds {
		if !c.IsAvailableCommand() {
			continue
		}

		emoji := c.Annotations[annotationEmoji]
		if emoji == "" {
			emoji = defaultEntryEmoji
		}

		group := c.Annotations[annotationGroup]
		if group == "" {
			group = GroupProject
		}

		entries = append(entries, Entry{
			Name:            c.Name(),
			Short:           c.Short,
			Long:            c.Long,
			Cmd:             c,
			Emoji:           emoji,
			RequiresProject: c.Annotations[annotationRequiresProject] == annotationTrue,
			RequiresDocker:  c.Annotations[annotationRequiresDocker] == annotationTrue,
			Group:           group,

			ShowsSetupProgress: c.Annotations[annotationShowsSetupProgress] == annotationTrue,
		})
	}

	return entries
}

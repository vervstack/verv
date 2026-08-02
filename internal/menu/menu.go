package menu

import (
	"github.com/spf13/cobra"
)

const (
	// annotationEmoji and annotationRequiresProject are the cobra.Command
	// Annotations keys commands set to opt into a picker icon and/or a
	// "needs an existing verv project" gate, keeping BuildEntries free of a
	// hardcoded command-name list.
	annotationEmoji           = "verv:emoji"
	annotationRequiresProject = "verv:requiresProject"

	defaultEntryEmoji = "▫️"
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

		entries = append(entries, Entry{
			Name:            c.Name(),
			Short:           c.Short,
			Long:            c.Long,
			Cmd:             c,
			Emoji:           emoji,
			RequiresProject: c.Annotations[annotationRequiresProject] == "true",
		})
	}

	return entries
}

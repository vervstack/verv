package actions

//go:generate minimock -i ActionPerformer -o ./../../../tests/mocks -g -s "_mock.go"

import (
	"fmt"
	"strings"
	"time"

	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/internal/io"
	"go.vervstack.ru/verv/internal/io/colors"
	"go.vervstack.ru/verv/plugins/project"
)

type Action interface {
	Do(p project.IProject) error
	NameInAction() string
}

type IActionPerformer interface {
	Tidy(proj project.IProject) error
}

type ActionPerformer struct {
	printer io.IO
}

func NewActionPerformer(printer io.IO) *ActionPerformer {
	return &ActionPerformer{
		printer: printer,
	}
}

func (a *ActionPerformer) Tidy(proj project.IProject) error {
	acts := GetTidyActionsForProject(proj.GetType())
	total := len(acts)

	a.printer.PrintlnColored(colors.ColorCyan, fmt.Sprintf("🚀 Tidying up %q — %d steps ahead", proj.GetName(), total))

	start := time.Now()

	for idx, ac := range acts {
		label := fmt.Sprintf("[%d/%d] %s %s", idx+1, total, stepEmoji(ac.NameInAction()), ac.NameInAction())

		spinner := io.NewSpinner(a.printer)
		spinner.Start(label)

		stepStart := time.Now()
		err := ac.Do(proj)
		elapsed := time.Since(stepStart).Round(time.Millisecond)

		if err != nil {
			spinner.Stop(false, fmt.Sprintf("%s — failed after %s", label, elapsed))

			return rerrors.Wrap(err)
		}

		spinner.Stop(true, fmt.Sprintf("%s — done in %s", label, elapsed))
	}

	totalElapsed := time.Since(start).Round(time.Millisecond)
	tidySummary := fmt.Sprintf("✨ Project is tidy! Wrapped up %d steps in %s", total, totalElapsed)
	a.printer.PrintlnColored(colors.ColorGreen, tidySummary)

	return nil
}

// stepEmojis maps a keyword found in a step's name to a small decorative icon.
// Purely cosmetic: an unmatched name falls back to a generic marker in stepEmoji.
var stepEmojis = []struct {
	keyword string
	emoji   string
}{
	{"workflow", "🔁"},
	{"commit", "📦"},
	{"hook", "🪝"},
	{"build", "🔨"},
	{"config", "⚙️"},
	{"clean", "🧹"},
	{"skeleton", "🏗️"},
	{"client", "🔌"},
	{"server", "🖥️"},
	{"dockerfile", "🐳"},
}

func stepEmoji(name string) string {
	n := strings.ToLower(name)

	for _, r := range stepEmojis {
		if strings.Contains(n, r.keyword) {
			return r.emoji
		}
	}

	return "▶️"
}

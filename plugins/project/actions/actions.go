package actions

//go:generate minimock -i IActionPerformer -n ActionPerformerMock -o ./../../../tests/mocks/action_performer_mock.go -g

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
	Tidy(proj project.IProject, fast, dirty bool) error
}

type ActionPerformer struct {
	printer io.IO
}

func NewActionPerformer(printer io.IO) *ActionPerformer {
	return &ActionPerformer{
		printer: printer,
	}
}

func (a *ActionPerformer) Tidy(proj project.IProject, fast, dirty bool) error {
	acts := GetTidyActionsForProject(proj.GetType(), fast, dirty)

	return RunPipeline(a.printer, proj, acts, PipelineLabels{
		StartEmoji: "🚀",
		StartVerb:  "Tidying up",
		EndEmoji:   "✨",
		EndVerb:    "is tidy",
	})
}

// PipelineLabels customizes the banner and completion messages printed by RunPipeline.
type PipelineLabels struct {
	StartEmoji string
	StartVerb  string
	EndEmoji   string
	EndVerb    string
}

// RunPipeline executes acts against proj, printing a colored banner, an animated
// spinner per step (à la `docker pull`), and a colored completion summary.
func RunPipeline(printer io.IO, proj project.IProject, acts []Action, labels PipelineLabels) error {
	total := len(acts)

	banner := fmt.Sprintf("%s %s %q — %d steps ahead", labels.StartEmoji, labels.StartVerb, proj.GetName(), total)
	printer.PrintlnColored(colors.ColorCyan, banner)

	start := time.Now()

	for idx, ac := range acts {
		label := fmt.Sprintf("[%d/%d] %s %s", idx+1, total, stepEmoji(ac.NameInAction()), ac.NameInAction())

		spinner := io.NewSpinner(printer)
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
	summary := fmt.Sprintf("%s Project %s! Wrapped up %d steps in %s",
		labels.EndEmoji, labels.EndVerb, total, totalElapsed)
	printer.PrintlnColored(colors.ColorGreen, summary)

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
	{"structure", "🧱"},
	{"init", "🌱"},
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

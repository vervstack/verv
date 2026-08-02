package menu

import (
	"os"
	"path/filepath"
	"strings"

	vervconfig "go.vervstack.ru/verv/internal/config"
	"go.vervstack.ru/verv/plugins/project"
)

const (
	goProjectEmoji    = "🐹"
	nonGoProjectEmoji = "📦"
	goModFileName     = "go.mod"
)

// Header carries the project-aware context shown above the interactive
// command picker.
type Header struct {
	IsVervProject bool
	Name          string
	Version       string
	Path          string
	Emoji         string
}

// BuildHeader inspects wd for the .verv/vervonomicon.yaml marker (the strict
// signal for "this is a verv project" — see project.IsVervProject) and
// assembles the picker header from it. A missing or unreadable version never
// fails the whole header; it's simply left blank.
func BuildHeader(wd string, cfg *vervconfig.VervConfig) Header {
	h := Header{
		Path: displayPath(wd),
	}

	h.IsVervProject = project.IsVervProject(wd)
	if !h.IsVervProject {
		return h
	}

	h.Emoji = projectEmoji(wd)

	vervonomicon, err := project.ReadVervonomicon(wd)
	if err != nil {
		h.IsVervProject = false

		return h
	}

	h.Name = vervonomicon.Name

	conf, err := project.LoadProjectConfig(wd, cfg)
	if err == nil {
		h.Version = conf.AppConfig.AppInfo.Version
	}

	return h
}

func projectEmoji(wd string) string {
	_, err := os.Stat(filepath.Join(wd, goModFileName))
	if err != nil {
		return nonGoProjectEmoji
	}

	return goProjectEmoji
}

// displayPath renders wd relative to the user's home directory (prefixed
// with ~) when possible, falling back to the raw wd on any error or when wd
// is outside the home directory.
func displayPath(wd string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return wd
	}

	rel, err := filepath.Rel(home, wd)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return wd
	}

	if rel == "." {
		return "~"
	}

	return filepath.Join("~", rel)
}

package init

import (
	"fmt"
	"path"
	"strings"

	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/internal/io"
	"go.vervstack.ru/verv/internal/io/colors"
	"go.vervstack.ru/verv/plugins/project/validators"
)

const (
	askUserForNameGreeting = "👋 Let's spin up something new!"
	askUserForNameQuestion = "   What should we call it?"
	askUserForNameHint     = `   💡 give it a full git path, like "github.com/vervstack/verv",
      or just a short name, like "verv" — we'll default to "%[1]s/verv"`
	askUserForNamePrompt = "➜  "

	nameAcceptedMessagePattern = `✨ Perfect — %q it is!`
)

var (
	errEmptyName = rerrors.New("no name entered")
)

type nameCollector struct {
	io io.IO

	defaultProjectGitPath string
}

func newNameCollector(io io.IO, defaultProjectGitPath string) nameCollector {
	return nameCollector{
		io:                    io,
		defaultProjectGitPath: defaultProjectGitPath,
	}
}

func (p *nameCollector) collect(args []string) (name string, err error) {
	if len(args) > 0 {
		name = args[0]
	} else {
		name, err = p.askUserForName()
		if err != nil {
			return "", rerrors.Wrap(err, "error while asking user for name")
		}
	}

	if name == "" {
		return "", rerrors.Wrap(errEmptyName)
	}

	name = p.removeHttpProtoc(name)
	name = p.preAppendHost(name)

	err = validators.ValidateProjectNameStr(name)
	if err != nil {
		return "", rerrors.Wrap(err, "error validating project name")
	}

	name = path.Join(path.Dir(name), strings.ToLower(path.Base(name)))

	p.io.PrintlnColored(colors.ColorCyan, fmt.Sprintf(nameAcceptedMessagePattern, name))

	return name, nil
}

func (p *nameCollector) askUserForName() (name string, err error) {
	p.io.PrintlnColored(colors.ColorMagenta, askUserForNameGreeting)
	p.io.Println()
	p.io.PrintlnColored(colors.ColorWhite, askUserForNameQuestion)
	p.io.PrintlnColored(colors.ColorYellow, fmt.Sprintf(askUserForNameHint, p.defaultProjectGitPath))
	p.io.Println()
	p.io.PrintColored(colors.ColorGreen, askUserForNamePrompt)

	name, err = p.io.GetInput()
	if err != nil {
		return "", rerrors.Wrap(err, "error obtaining project name")
	}

	return name, nil
}

func (p *nameCollector) removeHttpProtoc(name string) string {
	if strings.HasPrefix(name, "http") {
		return name[strings.Index(name, "://")+3:]
	}

	return name
}

func (p *nameCollector) preAppendHost(name string) string {
	firstDotIndex := strings.Index(name, ".")
	firstSlashIndex := strings.Index(name, "/")

	// if firstSlash comes after first dot - consider name already has host
	// TODO
	//nolint:mnd
	if firstSlashIndex-firstDotIndex > 2 {
		return name
	}

	return p.defaultProjectGitPath + "/" + name
}

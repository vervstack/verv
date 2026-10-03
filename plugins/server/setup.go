package server

import (
	"context"
	"os"

	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/internal/cmd"
	"go.vervstack.ru/verv/internal/io"
	"go.vervstack.ru/verv/internal/io/colors"
)

const (
	osReleasePath = "/etc/os-release"
	ubuntuId      = "ubuntu"
)

// Setup prepares the machine it runs on: deploy user, pubkey ssh, base
// packages, Docker and the user's authorized ssh keys. Ubuntu only, root only.
func Setup(printer io.IO, opts Options) error {
	err := checkPreconditions(opts)
	if err != nil {
		return rerrors.Wrap(err)
	}

	ctx := context.Background()
	summary := make([]string, 0, 6)

	steps := []step{
		{label: "Creating user " + opts.UserName, run: ensureUser},
		{label: "Setting up ssh", run: enablePubkeyAuthInSshd},
		{label: "Installing base packages", run: installBasePackages},
		{label: "Installing Docker", run: installDocker},
		{label: "Adding " + opts.UserName + " to docker and sudo groups", run: addUserToGroups},
		{label: "Installing verv for " + opts.UserName, run: installVervBinary},
	}

	if opts.SshKeyUrl == "" {
		printer.PrintlnColored(colors.ColorYellow,
			"variable SSH_KEY_URL is empty. You will need to install ssh key manually")
	} else {
		keysStep := step{label: "Installing ssh keys for " + opts.UserName, run: installSshKeys}
		steps = append(steps, keysStep)
	}

	for _, s := range steps {
		done, err := runStep(ctx, printer, s, opts)
		if err != nil {
			return rerrors.Wrap(err, "error during step: "+s.label)
		}

		summary = append(summary, done)
	}

	printSummary(printer, summary)

	return nil
}

type step struct {
	label string
	run   func(ctx context.Context, opts Options) (doneMessage string, err error)
}

func runStep(ctx context.Context, printer io.IO, s step, opts Options) (string, error) {
	spinner := io.NewSpinner(printer)
	spinner.Start(s.label)

	done, err := s.run(ctx, opts)
	if err != nil {
		spinner.Stop(false, s.label+" — failed")

		return "", rerrors.Wrap(err)
	}

	spinner.Stop(true, done)

	return done, nil
}

func printSummary(printer io.IO, summary []string) {
	printer.PrintlnColored(colors.ColorGreen, "Summary:")

	for _, line := range summary {
		printer.PrintlnColored(colors.ColorGreen, "* "+line)
	}
}

func checkPreconditions(opts Options) error {
	if os.Geteuid() != 0 {
		return rerrors.Wrap(errNotRoot)
	}

	content, err := os.ReadFile(osReleasePath)
	if err != nil {
		return rerrors.Wrap(err, "error reading os-release")
	}

	id, _ := parseOsRelease(string(content))
	if id != ubuntuId {
		return rerrors.Wrap(errUnsupportedOs)
	}

	if opts.UserName == "" {
		return rerrors.Wrap(errEmptyUserName)
	}

	if opts.Password == "" {
		return rerrors.Wrap(errEmptyPassword)
	}

	return nil
}

func run(tool string, args ...string) error {
	req := cmd.Request{Tool: tool, Args: args}

	_, err := cmd.Execute(req)
	if err != nil {
		return rerrors.Wrap(err)
	}

	return nil
}

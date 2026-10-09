package server

import (
	"context"
	"errors"
	"os"

	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/internal/cmd"
	"go.vervstack.ru/verv/internal/io"
	"go.vervstack.ru/verv/internal/io/colors"
	"go.vervstack.ru/verv/plugins/velez"
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
	steps := buildSteps(opts)
	summary := make([]string, 0, len(steps))

	if opts.SshKeyUrl == "" {
		printer.PrintlnColored(colors.ColorYellow,
			"variable SSH_KEY_URL is empty. You will need to install ssh key manually")
	}

	for _, s := range steps {
		if isStepSkippable(ctx, opts, s) {
			printer.PrintlnColored(colors.ColorGreen, "✅ "+s.doneLabel)

			summary = append(summary, s.doneLabel)

			continue
		}

		done, err := runStep(ctx, printer, s, opts)
		if err != nil {
			return rerrors.Wrap(err, "error during step: "+s.label)
		}

		summary = append(summary, done)
	}

	printSummary(printer, summary)

	return nil
}

// Progress returns how many setup steps are complete on this machine and how many exist.
// Read-only, no root needed.
func Progress(ctx context.Context, opts Options) (done, total int) {
	return countDone(ctx, opts, buildSteps(opts))
}

type step struct {
	label string
	// check reports the step's end state is already reached. Must be read-only and must work
	// without root (the menu calls it as a normal user).
	check func(ctx context.Context, opts Options) bool
	// forceRun, when non-nil and true, runs the step even if check passes.
	forceRun func(opts Options) bool
	// doneLabel is printed when the step is skipped because check passed.
	doneLabel string
	run       func(ctx context.Context, opts Options) (doneMessage string, err error)
}

func buildSteps(opts Options) []step {
	userName := opts.UserName

	steps := []step{
		{
			label:     "Creating user " + userName,
			check:     isUserSetUp,
			forceRun:  isPasswordChangeRequested,
			doneLabel: "User " + userName + " already exists",
			run:       ensureUser,
		},
		{
			label:     "Setting up ssh",
			check:     isSshdSetUp,
			doneLabel: "SSH already set up",
			run:       enablePubkeyAuthInSshd,
		},
		{
			label:     "Installing base packages",
			check:     areBasePackagesInstalled,
			doneLabel: "Base packages already installed",
			run:       installBasePackages,
		},
		{
			label:     "Installing vim",
			check:     isVimSetUp,
			doneLabel: "vim already installed and configured",
			run:       installVim,
		},
		{
			label:     "Installing Docker",
			check:     isDockerInstalled,
			doneLabel: "Docker already installed",
			run:       installDocker,
		},
		{
			label:     "Adding " + userName + " to docker and sudo groups",
			check:     isUserInDockerAndSudoGroups,
			doneLabel: userName + " already in docker and sudo groups",
			run:       addUserToGroups,
		},
		{
			label:     "Installing verv for " + userName,
			check:     isVervBinaryInstalled,
			doneLabel: "verv already installed to " + vervInstallPath,
			run:       installVervBinary,
		},
		{
			label:     "Preparing " + velez.DefaultKeyPath + " for Velez keys",
			check:     isVelezKeyDirSetUp,
			doneLabel: velez.DefaultKeyPath + " already prepared for Velez keys",
			run:       prepareVelezKeyDir,
		},
		{
			label:     "Installing sysbox",
			check:     isSysboxSetUp,
			doneLabel: "Sysbox already installed",
			run:       installSysbox,
		},
	}

	if opts.SshKeyUrl != "" {
		keysStep := step{
			label:     "Installing ssh keys for " + userName,
			check:     areSshKeysInstalled,
			forceRun:  isAlways,
			doneLabel: "SSH keys already installed for " + userName,
			run:       installSshKeys,
		}

		steps = append(steps, keysStep)
	}

	return steps
}

func isAlways(_ Options) bool {
	return true
}

func isStepSkippable(ctx context.Context, opts Options, s step) bool {
	if !s.check(ctx, opts) {
		return false
	}

	return s.forceRun == nil || !s.forceRun(opts)
}

func countDone(ctx context.Context, opts Options, steps []step) (done, total int) {
	for _, s := range steps {
		if s.check(ctx, opts) {
			done++
		}
	}

	return done, len(steps)
}

func runStep(ctx context.Context, printer io.IO, s step, opts Options) (string, error) {
	spinner := io.NewSpinner(printer)
	spinner.Start(s.label)

	done, err := s.run(ctx, opts)

	var skip *skipError

	if errors.As(err, &skip) {
		spinner.StopInfo(s.label + " — skipped")

		printer.PrintlnColored(colors.ColorYellow, "⚠ "+s.label+" skipped: "+skip.reason)

		return s.label + " skipped: " + skip.reason, nil
	}

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

	if opts.Password == "" && !UserExists(opts.UserName) {
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

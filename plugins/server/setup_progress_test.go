package server

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	testUserName = "deployer"
)

func newCheckStep(isDone bool) step {
	check := func(_ context.Context, _ Options) bool { return isDone }

	return step{check: check}
}

func Test_CountDone_Scenarios(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		steps     []step
		wantDone  int
		wantTotal int
	}{
		{"no steps", nil, 0, 0},
		{"none complete", []step{newCheckStep(false), newCheckStep(false)}, 0, 2},
		{"some complete", []step{newCheckStep(true), newCheckStep(false), newCheckStep(true)}, 2, 3},
		{"all complete", []step{newCheckStep(true), newCheckStep(true)}, 2, 2},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			done, total := countDone(context.Background(), Options{}, tc.steps)

			require.Equal(t, tc.wantDone, done)
			require.Equal(t, tc.wantTotal, total)
		})
	}
}

func Test_IsStepSkippable_Scenarios(t *testing.T) {
	t.Parallel()

	forceAlways := func(_ Options) bool { return true }
	forceNever := func(_ Options) bool { return false }

	cases := []struct {
		name string
		step step
		want bool
	}{
		{"check fails", step{check: newCheckStep(false).check}, false},
		{"check passes without forceRun", step{check: newCheckStep(true).check}, true},
		{"check passes and forceRun false", step{check: newCheckStep(true).check, forceRun: forceNever}, true},
		{"check passes and forceRun true", step{check: newCheckStep(true).check, forceRun: forceAlways}, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := isStepSkippable(context.Background(), Options{}, tc.step)

			require.Equal(t, tc.want, got)
		})
	}
}

func Test_BuildSteps_SshKeysStepOnlyWithUrl(t *testing.T) {
	t.Parallel()

	withoutUrl := buildSteps(Options{UserName: testUserName})
	withUrl := buildSteps(Options{UserName: testUserName, SshKeyUrl: "https://example.com/keys"})

	require.Len(t, withUrl, len(withoutUrl)+1)
}

func Test_IsPubkeyAuthEnabled_Scenarios(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		content string
		want    bool
	}{
		{"commented default", "Port 22\n#PubkeyAuthentication yes\n", false},
		{"enabled", "Port 22\nPubkeyAuthentication yes\n", true},
		{"unrelated config", "Port 22\n", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tc.want, isPubkeyAuthEnabled(tc.content))
		})
	}
}

func Test_IsDpkgInstalled_Scenarios(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		status string
		want   bool
	}{
		{"installed", "install ok installed", true},
		{"half configured", "install ok half-configured", false},
		{"removed with config", "deinstall ok config-files", false},
		{"empty", "", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tc.want, isDpkgInstalled(tc.status))
		})
	}
}

func Test_HasAllGroups_Scenarios(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		output string
		groups []string
		want   bool
	}{
		{"all present", "deployer sudo docker\n", []string{dockerGroup, sudoGroup}, true},
		{"one missing", "deployer sudo\n", []string{dockerGroup, sudoGroup}, false},
		{"substring is not a match", "deployer dockerroot\n", []string{dockerGroup}, false},
		{"empty output", "", []string{"verv"}, false},
		{"no groups required", testUserName, nil, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tc.want, hasAllGroups(tc.output, tc.groups...))
		})
	}
}

func Test_HasAuthorizedKeys_Scenarios(t *testing.T) {
	t.Parallel()

	require.False(t, hasAuthorizedKeys(" \n\n"))
	require.True(t, hasAuthorizedKeys("ssh-ed25519 AAAA key\n"))
}

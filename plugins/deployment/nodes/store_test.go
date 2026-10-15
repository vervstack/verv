package nodes

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func newTestStore(t *testing.T) (*Store, string, string) {
	t.Helper()

	home := t.TempDir()
	project := t.TempDir()
	t.Setenv("HOME", home)

	store, err := NewStore(project)
	require.NoError(t, err)

	return store, home, project
}

func enableProjectScope(t *testing.T, project string) {
	t.Helper()

	err := os.MkdirAll(filepath.Join(project, ".verv"), 0o755)
	require.NoError(t, err)
}

func newNode(name, url string) Node {
	return Node{Name: name, Url: url, Secret: "secret-" + name}
}

func Test_List_NoFilesReturnsEmpty(t *testing.T) {
	store, _, _ := newTestStore(t)

	entries, err := store.List()
	require.NoError(t, err)
	require.Empty(t, entries)
}

func Test_List_ProjectOverridesGlobalByName(t *testing.T) {
	store, _, project := newTestStore(t)
	enableProjectScope(t, project)

	err := store.Add(newNode("home", "https://global.example.com"), ScopeGlobal)
	require.NoError(t, err)
	err = store.Add(newNode("home", "https://project.example.com"), ScopeProject)
	require.NoError(t, err)

	entries, err := store.List()
	require.NoError(t, err)
	require.Len(t, entries, 1)
	require.Equal(t, "https://project.example.com", entries[0].Url)
	require.Equal(t, ScopeProject, entries[0].Scope)
}

func Test_List_SortedByNameWithScopes(t *testing.T) {
	store, _, project := newTestStore(t)
	enableProjectScope(t, project)

	err := store.Add(newNode("zeta", "https://z.example.com"), ScopeGlobal)
	require.NoError(t, err)
	err = store.Add(newNode("alpha", "https://a.example.com"), ScopeProject)
	require.NoError(t, err)
	err = store.Add(newNode("mid", "https://m.example.com"), ScopeGlobal)
	require.NoError(t, err)

	entries, err := store.List()
	require.NoError(t, err)
	require.Len(t, entries, 3)
	require.Equal(t, "alpha", entries[0].Name)
	require.Equal(t, ScopeProject, entries[0].Scope)
	require.Equal(t, "mid", entries[1].Name)
	require.Equal(t, ScopeGlobal, entries[1].Scope)
	require.Equal(t, "zeta", entries[2].Name)
}

func Test_Add_ReplacesSameName(t *testing.T) {
	store, _, _ := newTestStore(t)

	err := store.Add(newNode("home", "https://old.example.com"), ScopeGlobal)
	require.NoError(t, err)
	err = store.Add(newNode("home", "https://new.example.com"), ScopeGlobal)
	require.NoError(t, err)

	entries, err := store.List()
	require.NoError(t, err)
	require.Len(t, entries, 1)
	require.Equal(t, "https://new.example.com", entries[0].Url)
}

func Test_Add_WritesFileWithOwnerOnlyMode(t *testing.T) {
	store, home, project := newTestStore(t)
	enableProjectScope(t, project)

	err := store.Add(newNode("home", "https://a.example.com"), ScopeGlobal)
	require.NoError(t, err)
	err = store.Add(newNode("home", "https://b.example.com"), ScopeProject)
	require.NoError(t, err)

	globalInfo, err := os.Stat(filepath.Join(home, ".verv", "nodes.yaml"))
	require.NoError(t, err)
	require.Equal(t, os.FileMode(0o600), globalInfo.Mode().Perm())

	dirInfo, err := os.Stat(filepath.Join(home, ".verv"))
	require.NoError(t, err)
	require.Equal(t, os.FileMode(0o700), dirInfo.Mode().Perm())

	projectInfo, err := os.Stat(filepath.Join(project, ".verv", "nodes.local.yaml"))
	require.NoError(t, err)
	require.Equal(t, os.FileMode(0o600), projectInfo.Mode().Perm())
}

func Test_Add_ProjectScopeUnavailable(t *testing.T) {
	store, _, _ := newTestStore(t)

	require.False(t, store.IsProjectScopeAvailable())

	err := store.Add(newNode("home", "https://a.example.com"), ScopeProject)
	require.ErrorIs(t, err, errProjectScopeUnavailable)
}

func Test_Add_RejectsEmptyFields(t *testing.T) {
	store, _, _ := newTestStore(t)

	cases := []struct {
		name string
		node Node
		want error
	}{
		{"empty name", Node{Url: "u", Secret: "s"}, errEmptyName},
		{"empty url", Node{Name: "n", Secret: "s"}, errEmptyUrl},
		{"empty secret", Node{Name: "n", Url: "u"}, errEmptySecret},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := store.Add(tc.node, ScopeGlobal)
			require.ErrorIs(t, err, tc.want)
		})
	}
}

func Test_Add_ProjectScopeCreatesGitignore(t *testing.T) {
	store, _, project := newTestStore(t)
	enableProjectScope(t, project)

	err := store.Add(newNode("home", "https://a.example.com"), ScopeProject)
	require.NoError(t, err)

	content, err := os.ReadFile(filepath.Join(project, ".gitignore"))
	require.NoError(t, err)
	require.Equal(t, ".verv/nodes.local.yaml\n", string(content))
}

func Test_Add_ProjectScopeAppendsToGitignoreOnce(t *testing.T) {
	store, _, project := newTestStore(t)
	enableProjectScope(t, project)

	gitignorePath := filepath.Join(project, ".gitignore")
	err := os.WriteFile(gitignorePath, []byte("dist"), 0o644)
	require.NoError(t, err)

	err = store.Add(newNode("one", "https://a.example.com"), ScopeProject)
	require.NoError(t, err)
	err = store.Add(newNode("two", "https://b.example.com"), ScopeProject)
	require.NoError(t, err)

	content, err := os.ReadFile(gitignorePath)
	require.NoError(t, err)
	require.Equal(t, "dist\n.verv/nodes.local.yaml\n", string(content))
}

func Test_Add_ProjectScopeKeepsExistingGitignoreEntry(t *testing.T) {
	store, _, project := newTestStore(t)
	enableProjectScope(t, project)

	gitignorePath := filepath.Join(project, ".gitignore")
	original := ".verv/nodes.local.yaml\ndist\n"
	err := os.WriteFile(gitignorePath, []byte(original), 0o644)
	require.NoError(t, err)

	err = store.Add(newNode("home", "https://a.example.com"), ScopeProject)
	require.NoError(t, err)

	content, err := os.ReadFile(gitignorePath)
	require.NoError(t, err)
	require.Equal(t, original, string(content))
}

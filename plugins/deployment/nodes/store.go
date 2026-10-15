package nodes

import (
	"cmp"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"slices"

	"go.redsock.ru/rerrors"
	"gopkg.in/yaml.v3"
)

const (
	vervDirName         = ".verv"
	globalFileName      = "nodes.yaml"
	projectFileName     = "nodes.local.yaml"
	dirMode             = 0o700
	fileMode            = 0o600
	temporaryFilePrefix = ".nodes-*.tmp"
)

type Store struct {
	globalPath string
	projectDir string
}

func NewStore(projectDir string) (*Store, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, rerrors.Wrap(err, "error resolving home directory")
	}

	store := &Store{
		globalPath: filepath.Join(home, vervDirName, globalFileName),
		projectDir: projectDir,
	}

	return store, nil
}

func (s *Store) IsProjectScopeAvailable() bool {
	info, err := os.Stat(filepath.Join(s.projectDir, vervDirName))
	if err != nil {
		return false
	}

	return info.IsDir()
}

func (s *Store) List() ([]Entry, error) {
	globalNodes, err := readNodes(s.globalPath)
	if err != nil {
		return nil, rerrors.Wrap(err, "error reading global nodes")
	}

	projectNodes, err := readNodes(s.projectPath())
	if err != nil {
		return nil, rerrors.Wrap(err, "error reading project nodes")
	}

	byName := make(map[string]Entry, len(globalNodes)+len(projectNodes))
	for _, node := range globalNodes {
		byName[node.Name] = Entry{Node: node, Scope: ScopeGlobal}
	}
	for _, node := range projectNodes {
		byName[node.Name] = Entry{Node: node, Scope: ScopeProject}
	}

	entries := make([]Entry, 0, len(byName))
	for _, entry := range byName {
		entries = append(entries, entry)
	}
	slices.SortFunc(entries, compareEntries)

	return entries, nil
}

func (s *Store) Add(node Node, scope Scope) error {
	err := validateNode(node)
	if err != nil {
		return rerrors.Wrap(err, "error validating node")
	}

	path, err := s.pathForScope(scope)
	if err != nil {
		return rerrors.Wrap(err)
	}

	existing, err := readNodes(path)
	if err != nil {
		return rerrors.Wrap(err, "error reading nodes")
	}

	index := indexOfName(existing, node.Name)
	if index >= 0 {
		existing[index] = node
	} else {
		existing = append(existing, node)
	}

	err = writeNodes(path, existing)
	if err != nil {
		return rerrors.Wrap(err, "error writing nodes")
	}

	if scope == ScopeProject {
		err = ensureGitignored(s.projectDir)
		if err != nil {
			return rerrors.Wrap(err, "error ignoring project nodes file")
		}
	}

	return nil
}

func (s *Store) projectPath() string {
	return filepath.Join(s.projectDir, vervDirName, projectFileName)
}

func (s *Store) pathForScope(scope Scope) (string, error) {
	switch scope {
	case ScopeGlobal:
		return s.globalPath, nil
	case ScopeProject:
		if !s.IsProjectScopeAvailable() {
			return "", rerrors.Wrap(errProjectScopeUnavailable)
		}
		return s.projectPath(), nil
	default:
		return "", rerrors.Wrap(errUnknownScope)
	}
}

func validateNode(node Node) error {
	if node.Name == "" {
		return rerrors.Wrap(errEmptyName)
	}
	if node.Url == "" {
		return rerrors.Wrap(errEmptyUrl)
	}
	if node.Secret == "" {
		return rerrors.Wrap(errEmptySecret)
	}

	return nil
}

func indexOfName(nodes []Node, name string) int {
	for i, node := range nodes {
		if node.Name == name {
			return i
		}
	}

	return -1
}

func compareEntries(a, b Entry) int {
	return cmp.Compare(a.Name, b.Name)
}

func readNodes(path string) ([]Node, error) {
	raw, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, rerrors.Wrap(err, "error reading file")
	}

	var content fileContent
	err = yaml.Unmarshal(raw, &content)
	if err != nil {
		return nil, rerrors.Wrap(err, "error parsing yaml")
	}

	return content.Nodes, nil
}

func writeNodes(path string, nodes []Node) error {
	content := fileContent{Nodes: nodes}

	raw, err := yaml.Marshal(content)
	if err != nil {
		return rerrors.Wrap(err, "error marshaling yaml")
	}

	dir := filepath.Dir(path)
	err = os.MkdirAll(dir, dirMode)
	if err != nil {
		return rerrors.Wrap(err, "error creating directory")
	}

	return writeFileAtomic(dir, path, raw)
}

func writeFileAtomic(dir, path string, raw []byte) error {
	temporary, err := os.CreateTemp(dir, temporaryFilePrefix)
	if err != nil {
		return rerrors.Wrap(err, "error creating temporary file")
	}
	temporaryPath := temporary.Name()

	_, err = temporary.Write(raw)
	if err != nil {
		_ = temporary.Close()
		_ = os.Remove(temporaryPath)
		return rerrors.Wrap(err, "error writing temporary file")
	}

	err = temporary.Close()
	if err != nil {
		_ = os.Remove(temporaryPath)
		return rerrors.Wrap(err, "error closing temporary file")
	}

	err = os.Chmod(temporaryPath, fileMode)
	if err != nil {
		_ = os.Remove(temporaryPath)
		return rerrors.Wrap(err, "error setting file mode")
	}

	err = os.Rename(temporaryPath, path)
	if err != nil {
		_ = os.Remove(temporaryPath)
		return rerrors.Wrap(err, "error replacing file")
	}

	return nil
}

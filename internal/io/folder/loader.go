package folder

import (
	"os"
	"path"
	"strings"
)

func Load(root string, ops ...opt) (*Folder, error) {
	dir, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	o := opts{}
	for _, f := range ops {
		f(&o)
	}

	f := &Folder{
		Name:  root,
		Inner: make([]*Folder, 0, len(dir)),
	}
	for _, d := range dir {
		innerFolder, err := load(path.Join(root, d.Name()), "", o)
		if err != nil {
			return nil, err
		}
		if innerFolder != nil {
			f.Inner = append(f.Inner, innerFolder)
		}
	}

	return f, nil
}

func matchesAnyPattern(patterns []string, name, parent string) bool {
	relPath := path.Join(parent, name)

	for _, pattern := range patterns {
		pattern = strings.TrimRight(pattern, "/\r ")
		if pattern == "" {
			continue
		}

		if strings.Contains(pattern, "/") {
			if matched, _ := path.Match(pattern, relPath); matched {
				return true
			}
		} else {
			if matched, _ := path.Match(pattern, name); matched {
				return true
			}
		}
	}
	return false
}

func load(root, parent string, o opts) (*Folder, error) {
	if matchesAnyPattern(o.ignoredPaths, path.Base(root), parent) {
		return nil, nil
	}

	st, err := os.Stat(root)
	if err != nil {
		return nil, err
	}

	if !st.IsDir() {
		var innerFile []byte
		innerFile, err = os.ReadFile(root)
		if err != nil {
			return nil, err
		}

		folder := &Folder{
			Name:    path.Base(root),
			Content: innerFile,
		}

		folder.olderVersion = make([]byte, len(folder.Content))
		copy(folder.olderVersion, folder.Content)

		return folder, nil
	}

	dir, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	f := &Folder{
		Name:  path.Base(root),
		Inner: make([]*Folder, 0, len(dir)),
	}

	for _, d := range dir {
		var innerDir *Folder
		innerDir, err = load(path.Join(root, d.Name()), path.Join(parent, f.Name), o)
		if err != nil {
			return nil, err
		}

		if innerDir != nil {
			f.Inner = append(f.Inner, innerDir)
		}
	}

	return f, nil
}

package folder

import (
	"os"
	"path"
	"strings"

	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/internal/io"
)

type Folder struct {
	Name    string
	Inner   []*Folder
	Content []byte

	olderVersion  []byte
	isToBeDeleted bool
}

func (f *Folder) Add(folders ...*Folder) {
	for _, fl := range folders {
		dir := path.Dir(fl.Name)
		name := path.Base(fl.Name)

		if dir == "." {
			found := false

			for idx := range f.Inner {
				if f.Inner[idx].Name == fl.Name {
					f.Inner[idx] = fl
					found = true

					break
				}
			}

			if !found {
				f.Inner = append(f.Inner, fl)
			}
		} else {
			fl.Name = name
			f.addWithPath(dir, fl)
		}
	}
}

func (f *Folder) GetByPath(pth ...string) *Folder {
	currentFolder := f

	splitPath := make([]string, 0, len(pth))
	for _, p := range pth {
		sp := strings.Split(p, string(os.PathSeparator))

		splitPath = append(splitPath, sp...)
	}

	for _, p := range splitPath {
		var foundFolder *Folder

		for _, cf := range currentFolder.Inner {
			if cf.Name == p {
				foundFolder = cf

				break
			}
		}

		if foundFolder == nil {
			return nil
		}

		currentFolder = foundFolder
	}

	return currentFolder
}

func (f *Folder) Build() error {
	err := f.build(path.Dir(f.Name))
	if err != nil {
		return rerrors.Wrap(err, "folder build failed")
	}

	return nil
}

func (f *Folder) CopyWithNewName(name string) *Folder {
	newF := Folder{
		Name:    name,
		Content: make([]byte, len(f.Content)),
	}

	copy(newF.Content, f.Content)

	return &newF
}

func (f *Folder) Copy() *Folder {
	newF := Folder{
		Name:    f.Name,
		Content: make([]byte, len(f.Content)),
	}

	copy(newF.Content, f.Content)

	return &newF
}

func (f *Folder) Delete() {
	if f == nil {
		return
	}

	f.isToBeDeleted = true
}

func (f *Folder) findOrCreateChild(name string) *Folder {
	for idx := range f.Inner {
		if f.Inner[idx].Name == name {
			return f.Inner[idx]
		}
	}

	child := &Folder{Name: name}

	f.Inner = append(f.Inner, child)

	return child
}

func (f *Folder) mergeChild(folderToAdd *Folder) {
	for idx, itemInCurrentFolder := range f.Inner {
		if itemInCurrentFolder.Name != folderToAdd.Name {
			continue
		}

		if len(f.Inner[idx].Content) != 0 && len(folderToAdd.Content) != 0 {
			f.Inner[idx].Content = folderToAdd.Content
		} else {
			f.Inner[idx] = folderToAdd
		}

		return
	}

	f.Inner = append(f.Inner, folderToAdd)
}

func (f *Folder) addWithPath(pth string, folders ...*Folder) {
	if len(folders) == 0 {
		return
	}

	pths := strings.Split(pth, string(os.PathSeparator))

	currentFolder := f
	for _, pathPart := range pths {
		currentFolder = currentFolder.findOrCreateChild(pathPart)
	}

	for _, folderToAdd := range folders {
		currentFolder.mergeChild(folderToAdd)
	}
}

func (f *Folder) build(root string) error {
	pth := path.Join(root, path.Base(f.Name))

	if f.isToBeDeleted {
		return f.buildDelete(pth)
	}

	if len(f.Content) != 0 {
		err := f.buildFile(pth)
		if err != nil {
			return rerrors.Wrap(err, "failed to building file")
		}

		return nil
	}

	err := f.buildDir(pth)
	if err != nil {
		return rerrors.Wrap(err, "error building directory")
	}

	return nil
}

func (f *Folder) buildDelete(pth string) error {
	err := os.RemoveAll(pth)
	if err != nil {
		return rerrors.Wrap(err)
	}

	return nil
}

func (f *Folder) buildFile(pth string) error {
	if f.isUnchangedFromOlderVersion() {
		return nil
	}

	if len(f.Content) != 1 || f.Content[0] == 0 {
		err := io.OverrideFile(pth, f.Content)
		if err != nil {
			return rerrors.Wrap(err)
		}
	}

	return nil
}

func (f *Folder) isUnchangedFromOlderVersion() bool {
	if len(f.olderVersion) != len(f.Content) {
		return false
	}

	var idx int

	for idx = range f.olderVersion {
		if f.olderVersion[idx] != f.Content[idx] {
			break
		}
	}

	return len(f.olderVersion) != idx-1
}

func (f *Folder) buildDir(pth string) error {
	err := os.MkdirAll(pth, io.DefaultDirPerm)
	if err != nil {
		return rerrors.Wrap(err, "failed to create directory:", pth)
	}

	for _, d := range f.Inner {
		err = d.build(pth)
		if err != nil {
			return rerrors.Wrap(err, "failed to build")
		}
	}

	return nil
}

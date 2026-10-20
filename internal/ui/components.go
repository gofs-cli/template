package ui

import (
	"fmt"
	"io/fs"
	"path"
	"strings"
)

const elementFile = "element.js"

func (u *UI) ComponentFS() fs.FS {
	return componentFS{files: FS, paths: u.components}
}

func (u *UI) addComponent(p string) error {
	tag := tagFromPath(p)
	if tag == "" {
		return nil
	}
	if other, ok := u.components[tag]; ok {
		return fmt.Errorf("ui: component <%s> defined in both %s and %s", tag, other, p)
	}
	u.components[tag] = p
	return nil
}

func tagFromPath(p string) string {
	dir, file := path.Split(p)
	if file != elementFile {
		return ""
	}
	_, tag := path.Split(strings.TrimSuffix(dir, "/"))
	if !strings.Contains(tag, "-") {
		return ""
	}
	return tag
}

type componentFS struct {
	files fs.FS
	paths map[string]string
}

func (c componentFS) Open(name string) (fs.File, error) {
	tag, file, nested := strings.Cut(name, "/")
	if !nested {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
	}
	p, ok := c.paths[tag]
	if !ok || file != elementFile {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
	}
	return c.files.Open(p)
}

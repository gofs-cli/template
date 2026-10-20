package ui

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"path"
	"path/filepath"
	"sync"
)

//go:embed app
var FS embed.FS

const appFolder = "app"

type UI struct {
	root       *template.Template // every template (never executed so it can be cloned)
	fragments  *template.Template // clone of root used to execute fragments
	cache      sync.Map
	components map[string]string   // custom element tag -> path of its element.js
	templates  map[string]struct{} // unique template names
}

func Init() (*UI, error) {
	root := template.New("root").Funcs(funcs)
	if root == nil {
		return nil, fmt.Errorf("could not create root template")
	}
	u := &UI{root: root, components: map[string]string{}, templates: map[string]struct{}{}}
	err := fs.WalkDir(FS, appFolder, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(p) == ".js" {
			return u.addComponent(p)
		}
		if filepath.Ext(p) != ".gohtml" {
			return nil
		}
		t, err := template.New(path.Base(p)).Funcs(funcs).ParseFS(FS, p)
		if err != nil {
			return err
		}
		if t == nil {
			return nil // no templates in the file
		}
		for _, child := range t.Templates() {
			name := child.Name()
			if name == path.Base(p) || child.Tree == nil {
				continue // don't add file template or empty templates
			}
			if other, ok := u.templates[name]; ok {
				return fmt.Errorf("ui: template %q defined in both %s and %s", name, other, p)
			}
			u.templates[name] = struct{}{}
			if _, err := root.AddParseTree(name, child.Tree); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	u.fragments, err = root.Clone()
	if err != nil {
		return nil, err
	}
	return u, nil
}

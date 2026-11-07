// Package ui parses the templates embedded under app/ once, at startup, and
// renders them by name.
package ui

import (
	"bytes"
	"embed"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
)

//go:embed app
var appFS embed.FS

// templates holds every {{ define }} block found under app/. A page includes
// "layouts/head" and "layouts/foot" to form the document.
var templates = mustParse()

// mustParse reads every .gohtml file into one template set. Files are parsed
// one by one under their path name, so same-named files do not collide; a
// name defined twice keeps the last one seen, as html/template does.
func mustParse() *template.Template {
	t := template.New("root").Funcs(funcs)
	err := fs.WalkDir(appFS, "app", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || filepath.Ext(p) != ".gohtml" {
			return err
		}
		body, err := appFS.ReadFile(p)
		if err != nil {
			return err
		}
		_, err = t.New(p).Parse(string(body))
		return err
	})
	if err != nil {
		panic("ui: " + err.Error())
	}
	return t
}

// Render executes the named template. A missing or broken template is a
// programming error, so it answers 500 with the error text.
func Render(w http.ResponseWriter, name string, data any) {
	var buf bytes.Buffer
	if err := templates.ExecuteTemplate(&buf, name, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = buf.WriteTo(w)
}

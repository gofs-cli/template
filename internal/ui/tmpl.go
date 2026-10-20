package ui

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
)

func (u *UI) getCachedOrPrepare(layout string, content string) (*template.Template, error) {
	key := layout + content
	if tmpl, ok := u.cache.Load(key); ok {
		return tmpl.(*template.Template), nil
	}
	r, err := u.root.Clone()
	if err != nil {
		return nil, err
	}
	c := u.root.Lookup(content)
	if c == nil {
		return nil, fmt.Errorf("template not found: %s", content)
	}
	r.AddParseTree("content", c.Tree)
	u.cache.Store(key, r)
	return r, nil
}

func (u *UI) RenderLayoutContent(w http.ResponseWriter, layout string, content string, data any) {
	tmpl, err := u.getCachedOrPrepare(layout, content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	render(w, tmpl, layout, data)
}

func (u *UI) RenderFragment(w http.ResponseWriter, fragment string, data any) {
	if u.fragments.Lookup(fragment) == nil {
		http.Error(w, "template not found: "+fragment, http.StatusInternalServerError)
		return
	}
	render(w, u.fragments, fragment, data)
}

func (u *UI) WantsFragment(r *http.Request) bool {
	return r.Header.Get("Accept") == "text/html" && r.Header.Get("Sec-Fetch-Mode") != "navigate"
}

func render(w http.ResponseWriter, tmpl *template.Template, name string, data any) {
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, name, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = buf.WriteTo(w)
}

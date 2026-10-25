package ui

import (
	"fmt"
	"html/template"
)

// functions available in every template.
var funcs = template.FuncMap{
	"dict": dict,
}

// dict builds a map from key/value pairs so that a page can pass parameters
// to a shared template, such as the document head:
//
//	{{ template "layouts/head" dict "Title" "Contacts" }}
//
// Keys must be strings. Missing keys read as nil, so use {{ with .Key }} for
// optional parameters.
func dict(pairs ...any) (map[string]any, error) {
	if len(pairs)%2 != 0 {
		return nil, fmt.Errorf("dict: odd number of arguments")
	}
	m := make(map[string]any, len(pairs)/2)
	for i := 0; i < len(pairs); i += 2 {
		key, ok := pairs[i].(string)
		if !ok {
			return nil, fmt.Errorf("dict: key %v is not a string", pairs[i])
		}
		m[key] = pairs[i+1]
	}
	return m, nil
}

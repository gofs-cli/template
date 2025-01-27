package activesearch

import (
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/gofs-cli/template/internal/ui"
	"github.com/gofs-cli/template/internal/ui/components/header"
)

func Index() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(ui.IndexPage(layout(header.Header(), body()))).ServeHTTP(w, r)
	})
}

var rows = []Row{
	{
		FirstName: "Rick",
		LastName:  "Grimes",
		Email:     "rick.grimes@email.com",
	},
	{
		FirstName: "Thomas",
		LastName:  "Shelby",
		Email:     "thomas.shelby@email.com",
	},
	{
		FirstName: "Harvey",
		LastName:  "Specter",
		Email:     "harvey.specter@email.com",
	},
	{
		FirstName: "Rick",
		LastName:  "Astley",
		Email:     "rick.astley@email.com",
	},
	{
		FirstName: "Thomas",
		LastName:  "cook",
		Email:     "thomas.cook@email.com",
	},
}

func Search() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		searchQuery := strings.ToLower(r.PostFormValue("search"))
		filteredRows := []Row{}
		for _, v := range rows {
			if contains(v.FirstName, searchQuery) || contains(v.LastName, searchQuery) || contains(v.Email, searchQuery) {
				filteredRows = append(filteredRows, v)
			}
		}
		templ.Handler(searchResults(filteredRows)).ServeHTTP(w, r)
	})
}

func contains(str, substr string) bool {
	return strings.Contains(strings.ToLower(str), substr)
}

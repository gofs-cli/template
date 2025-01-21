package clicktoedit

import (
	"net/http"

	"github.com/a-h/templ"

	"github.com/gofs-cli/template/internal/ui"
	"github.com/gofs-cli/template/internal/ui/components/header"
)

func Index() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(ui.IndexPage(layout(header.Header(), body()))).ServeHTTP(w, r)
	})
}

func Form() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(form("Joe", "Blow", "joe@blow.com")).ServeHTTP(w, r)
	})
}

func SaveForm() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fristName := r.PostFormValue("firstName")
		lastName := r.PostFormValue("lastName")
		email := r.PostFormValue("email")
		templ.Handler(form(fristName, lastName, email)).ServeHTTP(w, r)
	})
}

func EditForm() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(editForm()).ServeHTTP(w, r)
	})
}

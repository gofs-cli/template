package inlinevalidation

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

func Validate() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email := r.PostFormValue("email")
		if email != "test@test.com" {
			templ.Handler(errorEmail(email)).ServeHTTP(w, r)
		} else {
			templ.Handler(validEmail(email)).ServeHTTP(w, r)
		}
	})
}

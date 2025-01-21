package bulkupdate

import (
	"fmt"
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
		templ.Handler(form()).ServeHTTP(w, r)
	})
}

func Update() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var activated int
		var deactivated int

		for k, v := range initialState {
			a := r.PostFormValue("active:"+k) == "on"
			if a && !v {
				initialState[k] = true
				activated++
			} else if !a && v {
				initialState[k] = false
				deactivated++
			}
		}
		fmt.Fprintf(w, "Activated %d and deactivated %d users", activated, deactivated)
	})
}

package animations

import (
	"math/rand/v2"
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

func Colors() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		colors := []string{
			"red",
			"blue",
			"green",
			"orange",
			"pink",
		}

		// Fix: Ensure we use all colors
		color := colors[rand.IntN(len(colors))]

		// Fix: Pass color correctly to the templ component
		templ.Handler(newColor(color)).ServeHTTP(w, r)
	})
}

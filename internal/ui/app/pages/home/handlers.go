package home

import (
	"net/http"

	"github.com/gofs-cli/template/internal/ui"
)

func Index() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ui.Render(w, "home/index", nil)
	})
}

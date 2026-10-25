package notfound

import (
	"net/http"

	"github.com/gofs-cli/template/internal/ui"
)

func Index() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		ui.Render(w, "notfound/index", nil)
	})
}

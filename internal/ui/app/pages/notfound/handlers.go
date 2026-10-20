package notfound

import (
	"net/http"

	"github.com/gofs-cli/template/internal/ui"
)

func Index(ui *ui.UI) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		ui.RenderLayoutContent(w, "layouts/main", "notfound/index", nil)
	})
}

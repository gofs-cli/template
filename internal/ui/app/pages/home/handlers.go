package home

import (
	"net/http"

	"github.com/gofs-cli/template/internal/ui"
)

func Index(ui *ui.UI) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ui.RenderLayoutContent(w, "layouts/main", "home/index", nil)
	})
}

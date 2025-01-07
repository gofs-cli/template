package page

import (
	"net/http"

	"github.com/gofs-cli/template/internal/ui"

	"github.com/a-h/templ"
)

func Index() http.Handler {
	return templ.Handler(ui.Index())
}

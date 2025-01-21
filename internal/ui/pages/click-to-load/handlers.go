package clicktoload

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/a-h/templ"

	"github.com/gofs-cli/template/internal/ui"
	"github.com/gofs-cli/template/internal/ui/components/header"
)

func Index() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(ui.IndexPage(layout(header.Header(), body()))).ServeHTTP(w, r)
	})
}

func Page() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := r.URL.Query().Get("page")
		pInt, _ := strconv.Atoi(p)
		fmt.Printf("Page: %d\n", pInt)
		templ.Handler(page(pInt)).ServeHTTP(w, r)
	})
}

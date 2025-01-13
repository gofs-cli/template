package home

import (
	"net/http"

	"github.com/a-h/templ"

	"github.com/gofs-cli/template/internal/ui"
	"github.com/gofs-cli/template/internal/ui/components/header"
	"github.com/gofs-cli/template/internal/ui/components/toast"
)

func Index() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(ui.IndexPage(layout(header.Header(), body()))).ServeHTTP(w, r)
	})
}

func Success() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		toast.Success(w, r, "Success!")
	})
}

func Info() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		toast.Info(w, r, "Info!")
	})
}

func Warning() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		toast.Warning(w, r, "Warning!")
	})
}

func Error() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		toast.Error(w, r, http.StatusInternalServerError, "Error!")
	})
}

func Modal() http.Handler {
	return templ.Handler(ModalDemo())
}

package server

import (
	"net/http"

	"github.com/gofs-cli/template/internal/server/assets"
	"github.com/gofs-cli/template/internal/server/handlers"
	"github.com/gofs-cli/template/internal/ui/pages/home"
	"github.com/gofs-cli/template/internal/ui/pages/page1"
	"github.com/gofs-cli/template/internal/ui/pages/page2"
	"github.com/gofs-cli/template/internal/ui/pages/validation"
)

func (s *Server) Routes() {
	// filserver route for assets
	assetMux := http.NewServeMux()
	assetMux.Handle("GET /{path...}", http.StripPrefix("/assets/", handlers.NewHashedAssets(assets.FS)))
	s.r.Handle("GET /assets/{path...}", s.assetsMiddlewares(assetMux))

	// handlers for normal routes with all general middleware
	routesMux := http.NewServeMux()
	routesMux.Handle("GET /{$}", home.Index())
	routesMux.Handle("GET /modal", home.Modal())

	routesMux.Handle("GET /page1", page1.Index())
	routesMux.Handle("GET /page2", page2.Index())

	routesMux.Handle("/validate", validation.HandleNameValidation())

	routesMux.Handle("GET /toast-success", home.Success())
	routesMux.Handle("GET /toast-info", home.Info())
	routesMux.Handle("GET /toast-warning", home.Warning())
	routesMux.Handle("GET /toast-error", home.Error())

	s.r.Handle("/", s.routeMiddlewares(routesMux))

	s.srv.Handler = s.r
}

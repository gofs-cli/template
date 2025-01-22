package server

import (
	"net/http"

	"github.com/gofs-cli/template/internal/server/assets"
	"github.com/gofs-cli/template/internal/server/handlers"
	bulkupdate "github.com/gofs-cli/template/internal/ui/pages/bulk-update"
	clicktoedit "github.com/gofs-cli/template/internal/ui/pages/click-to-edit"
	clicktoload "github.com/gofs-cli/template/internal/ui/pages/click-to-load"
	"github.com/gofs-cli/template/internal/ui/pages/home"
	inlinevalidation "github.com/gofs-cli/template/internal/ui/pages/inline-validation"
	"github.com/gofs-cli/template/internal/ui/pages/notfound"
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
	routesMux.Handle("GET /", notfound.Index())

	// click to edit example
	routesMux.Handle("GET /click-to-edit", clicktoedit.Index())
	routesMux.Handle("GET /click-to-edit/contact/1/edit", clicktoedit.EditForm())
	routesMux.Handle("GET /click-to-edit/contact/1", clicktoedit.Form())
	routesMux.Handle("PUT /click-to-edit/contact/1", clicktoedit.SaveForm())

	// bulk update example
	routesMux.Handle("GET /bulk-update", bulkupdate.Index())
	routesMux.Handle("POST /bulk-update/users", bulkupdate.Update())

	// click to load example
	routesMux.Handle("GET /click-to-load", clicktoload.Index())
	routesMux.Handle("GET /click-to-load/contacts", clicktoload.Page())

	// inline validation example
	routesMux.Handle("GET /inline-validation", inlinevalidation.Index())
	routesMux.Handle("POST /inline-validation/email", inlinevalidation.Validate())

	routesMux.Handle("GET /modal", home.Modal())

	routesMux.Handle("GET /page1", page1.Index())
	routesMux.Handle("GET /page2", page2.Index())

	routesMux.Handle("POST /api/validate", validation.HandleNameValidation())

	routesMux.Handle("GET /toast-success", home.Success())
	routesMux.Handle("GET /toast-info", home.Info())
	routesMux.Handle("GET /toast-warning", home.Warning())
	routesMux.Handle("GET /toast-error", home.Error())

	s.r.Handle("/", s.routeMiddlewares(routesMux))

	s.srv.Handler = s.r
}

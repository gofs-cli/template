package server

import (
	"net/http"

	"github.com/gofs-cli/template/internal/server/assets"
	"github.com/gofs-cli/template/internal/server/handlers"
	"github.com/gofs-cli/template/internal/ui/app/pages/contacts"
	"github.com/gofs-cli/template/internal/ui/app/pages/home"
	"github.com/gofs-cli/template/internal/ui/app/pages/notfound"
)

func (s *Server) Routes() {
	// filserver route for assets
	assetMux := http.NewServeMux()
	assetMux.Handle("GET /{path...}", http.StripPrefix("/assets/", handlers.NewHashedAssets(assets.FS)))
	assetMux.Handle("GET /assets/components/{path...}", http.StripPrefix("/assets/components/", handlers.NewHashedAssets(s.ui.ComponentFS())))
	s.r.Handle("GET /assets/{path...}", s.assetsMiddlewares(assetMux))

	// handlers for normal routes with all general middleware
	routesMux := http.NewServeMux()
	routesMux.Handle("GET /{$}", home.Index(s.ui))
	routesMux.Handle("GET /", notfound.Index(s.ui))

	routesMux.Handle("GET /contacts", contacts.Index(s.ui, s.contacts))
	routesMux.Handle("GET /contacts/new", contacts.New(s.ui))
	routesMux.Handle("POST /contacts", contacts.Create(s.ui, s.contacts))
	routesMux.Handle("POST /contacts/validate", contacts.Validate(s.ui, s.contacts))
	routesMux.Handle("GET /contacts/{id}", contacts.Row(s.ui, s.contacts))
	routesMux.Handle("GET /contacts/{id}/edit", contacts.EditRow(s.ui, s.contacts))
	routesMux.Handle("POST /contacts/{id}", contacts.Update(s.ui, s.contacts))
	routesMux.Handle("POST /contacts/{id}/delete", contacts.Delete(s.ui, s.contacts))

	s.r.Handle("/", s.routeMiddlewares(routesMux))

	s.srv.Handler = s.r
}

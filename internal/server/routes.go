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
	// file server route for assets
	assetMux := http.NewServeMux()
	assetMux.Handle("GET /{path...}", http.StripPrefix("/assets/", handlers.NewHashedAssets(assets.FS)))
	s.r.Handle("GET /assets/{path...}", s.assetsMiddlewares(assetMux))

	// handlers for normal routes with all general middleware
	routesMux := http.NewServeMux()
	routesMux.Handle("GET /{$}", home.Index())
	routesMux.Handle("GET /", notfound.Index())

	routesMux.Handle("GET /contacts", contacts.Index(s.contacts))
	routesMux.Handle("GET /contacts/new", contacts.New())
	routesMux.Handle("POST /contacts", contacts.Create(s.contacts))
	routesMux.Handle("POST /contacts/{id}", contacts.Update(s.contacts))
	routesMux.Handle("POST /contacts/{id}/delete", contacts.Delete(s.contacts))

	s.r.Handle("/", s.routeMiddlewares(routesMux))

	s.srv.Handler = s.r
}

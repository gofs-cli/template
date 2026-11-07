package server

import (
	"net/http"
	"slices"

	"github.com/gofs-cli/template/internal/auth"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (s *Server) assetsMiddlewares(h http.Handler) http.Handler {
	middlewares := []func(http.Handler) http.Handler{
		cors.Handler(cors.Options{
			AllowedOrigins:   s.conf.AllowedOrigins,
			AllowedMethods:   []string{http.MethodGet, http.MethodOptions},
			AllowedHeaders:   []string{"X-CSRF-Token"},
			AllowCredentials: true,
			MaxAge:           300,
		}),
		middleware.Logger,
		middleware.Compress(5),
		// Asset URLs are not versioned, so revalidate instead of trusting a
		// long max-age: after a deploy the browser would otherwise keep stale
		// component JS (and its fragment URLs) for a week. The ETag makes a
		// revalidation a 304.
		middleware.SetHeader("Cache-Control", "no-cache"),
	}
	slices.Reverse(middlewares)
	for _, m := range middlewares {
		h = m(h)
	}
	return h
}

func (s *Server) routeMiddlewares(h http.Handler) http.Handler {
	middlewares := []func(http.Handler) http.Handler{
		cors.Handler(cors.Options{
			AllowedOrigins:   s.conf.AllowedOrigins,
			AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions},
			AllowedHeaders:   []string{"Authorization", "X-CSRF-Token"},
			AllowCredentials: true,
			MaxAge:           300,
		}),
		middleware.Logger,
		middleware.NoCache,
		middleware.StripSlashes,
		middleware.Recoverer,
		middleware.Compress(5),
		auth.Middleware(s.conf.Env),
	}
	slices.Reverse(middlewares)
	for _, m := range middlewares {
		h = m(h)
	}
	return h
}

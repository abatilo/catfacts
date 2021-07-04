package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

func (s *Server) registerRoutes() {
	s.router.Route("/", func(r chi.Router) {
		r.Get("/ping", s.ping())
	})
}

func (s *Server) ping() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pong")
	}
}

package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/go-chi/chi"
)

func (s *Server) registerRoutes() {
	s.router.Route("/api", func(r chi.Router) {
		r.Post("/sms/receive", s.receive())
		r.Get("/ping", s.ping())
	})
}

func (s *Server) ping() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pong")
	}
}

func (s *Server) receive() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			s.logger.Err(err).Msg("Couldn't read body")
			http.Error(w, "Couldn't read body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		m, _ := url.ParseQuery(string(body))
		s.logger.Info().Str("body", string(body)).Str("sms_body", m["Body"][0]).Msg("Received")

		fmt.Fprintf(w, "")
	}
}

package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/abatilo/catfacts/internal/model"
	"github.com/go-chi/chi"
	openapi "github.com/twilio/twilio-go/rest/lookups/v1"
	"gorm.io/gorm"
)

func (s *Server) registerRoutes() {
	s.router.Route("/api", func(r chi.Router) {
		r.Post("/sms/receive", s.receive())
		r.Get("/ping", s.ping())

		r.Post("/register", s.register())
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

func (s *Server) register() http.HandlerFunc {

	type registerRequest struct {
		PhoneNumber string
	}

	type registerResponse struct {
		PhoneNumber string `json:"phoneNumber,omitempty"`
		Active      bool   `json:"active"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// Parse request JSON
		var req registerRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			s.logger.Err(err).Msg("Bad request format")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		io.CopyN(ioutil.Discard, r.Body, 512)
		r.Body.Close()

		// Sanitize phone number
		countryCode := "US"
		fetchPhoneNumberResponse, err := s.twilioClient.LookupsV1.FetchPhoneNumber(req.PhoneNumber, &openapi.FetchPhoneNumberParams{
			CountryCode: &countryCode,
		})

		if err != nil {
			s.logger.Err(err).Msg("Couldn't look up this phone number")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		sanitized := *fetchPhoneNumberResponse.PhoneNumber

		// Place into database if it doesn't already exist
		target := model.Target{PhoneNumber: sanitized}
		result := s.db.Where(&target, "PhoneNumber").First(&target)

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			s.logger.Info().Str("phoneNumber", sanitized).Msg("Phone number wasn't found in DB")
			s.db.Create(&target)
		}

		json.NewEncoder(w).Encode(&registerResponse{
			PhoneNumber: target.PhoneNumber,
			Active:      target.Active,
		})
	}
}

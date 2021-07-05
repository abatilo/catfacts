package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/abatilo/catfacts/internal/facts"
	"github.com/abatilo/catfacts/internal/model"
	"github.com/go-chi/chi"
	tw_api "github.com/twilio/twilio-go/rest/api/v2010"
	tw_lookups "github.com/twilio/twilio-go/rest/lookups/v1"
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
		from := m["From"][0]
		smsBody := m["Body"][0]

		s.logger.Info().Str("body", smsBody).Str("from", from).Msg("Received")
		// Dispatch to commands
		switch strings.ToLower(smsBody) {
		case "y":
			target := model.Target{PhoneNumber: from}
			result := s.db.Where(&target, "PhoneNumber").First(&target)

			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				s.logger.Info().Str("phoneNumber", from).Msg("Phone number wasn't found in DB, creating now")
				s.db.Create(&target)
			}

			target.Active = true
			s.db.Save(&target)

			msg := "You've just been confirmed for Aaron Batilo's CatFacts! You will start receiving random CatFacts. You can text \"now\" if you'd like to immediately receive a CatFact"
			s.twilioClient.ApiV2010.CreateMessage(&tw_api.CreateMessageParams{
				From: &s.config.TwilioPhoneNumber,
				To:   &from,
				Body: &msg,
			})

			randomFact := facts.RandomFact()
			s.twilioClient.ApiV2010.CreateMessage(&tw_api.CreateMessageParams{
				From: &s.config.TwilioPhoneNumber,
				To:   &from,
				Body: &randomFact,
			})
		case "now":
			target := model.Target{PhoneNumber: from}
			s.db.Where(&target, "PhoneNumber").First(&target)

			if target.Active {
				randomFact := facts.RandomFact()
				s.twilioClient.ApiV2010.CreateMessage(&tw_api.CreateMessageParams{
					From: &s.config.TwilioPhoneNumber,
					To:   &from,
					Body: &randomFact,
				})
			} else {
				msg := "It doesn't look like this number has subscribed to CatFacts. Visit https://catfacts.aaronbatilo.dev if you'd like to change that!"
				s.twilioClient.ApiV2010.CreateMessage(&tw_api.CreateMessageParams{
					From: &s.config.TwilioPhoneNumber,
					To:   &from,
					Body: &msg,
				})
			}
		}

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
		fetchPhoneNumberResponse, err := s.twilioClient.LookupsV1.FetchPhoneNumber(req.PhoneNumber, &tw_lookups.FetchPhoneNumberParams{
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
			s.logger.Info().Str("phoneNumber", sanitized).Msg("Phone number wasn't found in DB, creating now")
			s.db.Create(&target)
		}

		// Send confirmation text
		if !target.Active {
			msg := "You've just been registered for Aaron Batilo's CatFacts! Reply with \"Y\" if you'd like to confirm that you want to receive CatFacts!"
			s.twilioClient.ApiV2010.CreateMessage(&tw_api.CreateMessageParams{
				From: &s.config.TwilioPhoneNumber,
				To:   &sanitized,
				Body: &msg,
			})
		}

		fmt.Fprintf(w, "")
	}
}

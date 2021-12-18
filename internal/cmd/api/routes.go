package api

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/abatilo/catfacts/internal/facts"
	"github.com/abatilo/catfacts/internal/model"
	"github.com/go-chi/chi"
	tw_api "github.com/twilio/twilio-go/rest/api/v2010"
	tw_lookups "github.com/twilio/twilio-go/rest/lookups/v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func (s *Server) registerRoutes() {
	s.router.Route("/api", func(r chi.Router) {
		r.Post("/sms/receive", s.receive())
		r.Get("/ping", s.ping())

		r.Post("/register", s.register())
	})
}

func (s *Server) connectToDB() (*gorm.DB, func() error) {
	s.logger.Info().Msg("Lazily instantiating a database connection")
	db, err := gorm.Open(postgres.Open(s.dbConnString), &gorm.Config{})
	if err != nil {
		s.logger.Panic().Err(err).Msg("Unable to connect to database")
	}

	raw, _ := db.DB()

	s.logger.Info().Msg("Starting migrations")
	db.AutoMigrate(
		&model.Target{},
	)
	s.logger.Info().Msg("Finished migrations")
	return db, raw.Close
}

func (s *Server) ping() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pong")
	}
}

func (s *Server) receive() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info().Msg("Received SMS")
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			s.logger.Err(err).Msg("Couldn't read body")
			http.Error(w, "Couldn't read body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		signatureString := s.config.TwilioHost + r.URL.String()

		postForm, _ := url.ParseQuery(string(body))
		keys := make([]string, 0, len(postForm))

		for key := range postForm {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		for _, key := range keys {
			signatureString += key + postForm[key][0]
		}

		mac := hmac.New(sha1.New, []byte(s.config.TwilioAuthToken))
		mac.Write([]byte(signatureString))
		expectedMac := mac.Sum(nil)
		expectedTwilioSignature := base64.StdEncoding.EncodeToString(expectedMac)

		if expectedTwilioSignature != r.Header.Get("X-Twilio-Signature") {
			s.logger.Info().Msg("Received request that didn't come from Twilio")
			http.Error(w, "Couldn't verify that the request came from Twilio", http.StatusUnauthorized)
			return
		}

		db, disconnect := s.connectToDB()
		defer disconnect()

		from := postForm["From"][0]
		smsBody := postForm["Body"][0]

		// Dispatch to commands
		switch strings.ToLower(smsBody) {
		case "y":
			target := model.Target{PhoneNumber: from}
			result := db.Where("phone_number = ?", from).First(&target)

			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				s.logger.Info().Str("phoneNumber", from).Msg("Phone number wasn't found in DB, creating now")
				db.Create(&target)
			}

			if !target.Active {
				go func() {
					db, disconnect := s.connectToDB()
					defer disconnect()

					msg := "You've just been confirmed for Aaron Batilo's CatFacts! You will start receiving random CatFacts. You can text \"now\" if you'd like to immediately receive a CatFact"
					s.twilioClient.ApiV2010.CreateMessage(&tw_api.CreateMessageParams{
						From: &s.config.TwilioPhoneNumber,
						To:   &from,
						Body: &msg,
					})

					randomFact, _ := facts.GenerateFact(target.ID)
					s.twilioClient.ApiV2010.CreateMessage(&tw_api.CreateMessageParams{
						From: &s.config.TwilioPhoneNumber,
						To:   &from,
						Body: &randomFact,
					})
					target.Active = true
					target.LastSMS = time.Now().UTC()
					db.Save(&target)
				}()
			} else {
				s.logger.Info().Str("phoneNumber", target.PhoneNumber).Msg("Phone number just tried to subscribe again")
			}

		case "now":
			target := model.Target{PhoneNumber: from}
			db.Where(&target, "PhoneNumber").First(&target)

			if target.Active {
				s.logger.Info().Msg("Calling goroutine")
				go func() {
					s.logger.Info().Msg("Starting goroutine")
					defer s.logger.Info().Msg("Completed goroutine")
					db, disconnect := s.connectToDB()
					defer disconnect()

					randomFact, _ := facts.GenerateFact(target.ID)
					s.twilioClient.ApiV2010.CreateMessage(&tw_api.CreateMessageParams{
						From: &s.config.TwilioPhoneNumber,
						To:   &from,
						Body: &randomFact,
					})

					target.LastSMS = time.Now().UTC()
					db.Save(&target)
				}()
			} else {
				msg := "It doesn't look like this number has subscribed to CatFacts. Visit https://catfacts.aaronbatilo.dev if you'd like to change that!"
				s.twilioClient.ApiV2010.CreateMessage(&tw_api.CreateMessageParams{
					From: &s.config.TwilioPhoneNumber,
					To:   &from,
					Body: &msg,
				})
			}
		}

		s.logger.Info().Msg("Writing to response")
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

		db, disconnect := s.connectToDB()
		defer disconnect()

		// Place into database if it doesn't already exist
		target := model.Target{PhoneNumber: sanitized}
		result := db.Where(&target, "PhoneNumber").First(&target)

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			s.logger.Info().Str("phoneNumber", sanitized).Msg("Phone number wasn't found in DB, creating now")
			db.Create(&target)
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

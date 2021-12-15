package api

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/pprof"

	gosundheit "github.com/AppsFlyer/go-sundheit"
	healthhttp "github.com/AppsFlyer/go-sundheit/http"
	"github.com/go-chi/chi"
	"github.com/twilio/twilio-go"
	"gorm.io/gorm"

	"github.com/rs/cors"
	"github.com/rs/zerolog"
)

const (
	// FlagPortName is the name for the flag that's used for serving the application
	FlagPortName = "PORT"

	// FlagPortDefault is the default value for the application web server
	FlagPortDefault = 8080

	// FlagAdminPortName is the name for the flag that's used for serving the application's administrative endpoints
	FlagAdminPortName = "ADMIN_PORT"

	// FlagAdminPortDefault is the default value for the application web server's administrative port
	FlagAdminPortDefault = 8081

	// FlagTwilioHostName is the flag for setting the Twilio host that's used for authenticating webhooks
	FlagTwilioHostName = "TWILIO_HOST"

	// FlagTwilioHostDefault is the default value of the TWILIO_HOST flag
	FlagTwilioHostDefault = ""

	// FlagTwilioAccountSIDName is the name of the flag for the configured Twilio Account String ID
	FlagTwilioAccountSIDName = "TWILIO_ACCOUNT_SID"

	// FlagTwilioAccountSIDDefault is the default value of the TWILIO_ACCOUNT_SID flag
	FlagTwilioAccountSIDDefault = ""

	// FlagTwilioAuthTokenName is the name of the flag for the configured Twilio Auth Token
	FlagTwilioAuthTokenName = "TWILIO_AUTH_TOKEN"

	// FlagTwilioAuthTokenDefault is the default for TWILIO_AUTH_TOKEN
	FlagTwilioAuthTokenDefault = "TWILIO_AUTH_TOKEN"

	// FlagTwilioPhoneNumberName is the flag for the FROM number
	FlagTwilioPhoneNumberName = "TWILIO_PHONE_NUMBER"

	// FlagTwilioPhoneNumberDefault is the default value of the TWILIO_PHONE_NUMBER flag
	FlagTwilioPhoneNumberDefault = ""

	FlagDBHost        = "DB_HOST"
	FlagDBHostDefault = "postgresql"

	FlagDBUser        = "DB_USER"
	FlagDBUserDefault = "postgres"

	FlagDBPassword        = "DB_PASSWORD"
	FlagDBPasswordDefault = "local_password"

	FlagDBName        = "DB_NAME"
	FlagDBNameDefault = "postgres"

	FlagDBSSLMode        = "DB_SSL_MODE"
	FlagDBSSLModeDefault = "disable"

	FlagDBSearchPath        = "DB_SEARCH_PATH"
	FlagDBSearchPathDefault = "public"
)

// Config is all configuration for running the application.
//
// We use a config struct so that we can statically type and check configuration values
type Config struct {
	// Port is the HTTP server port
	Port int

	// AdminPort is the HTTP server port for internal use
	AdminPort int

	// Twilio values
	TwilioHost        string
	TwilioAccountSID  string
	TwilioAuthToken   string
	TwilioPhoneNumber string

	DBHost       string
	DBUser       string
	DBPassword   string
	DBName       string
	DBSSLMode    string
	DBSearchPath string
}

// Server represents the service itself and all of its dependencies.
//
// This pattern is heavily based on the following blog post:
// https://pace.dev/blog/2018/05/09/how-I-write-http-services-after-eight-years.html
type Server struct {
	adminServer  *http.Server
	config       *Config
	logger       zerolog.Logger
	router       *chi.Mux
	server       *http.Server
	twilioClient *twilio.RestClient
	db           *gorm.DB
}

// ServerOption lets you functionally control construction of the web server
type ServerOption func(s *Server)

// NewServer creates a new api server
func NewServer(cfg *Config, options ...ServerOption) *Server {
	router := chi.NewRouter()
	s := &Server{
		config: cfg,
		logger: zerolog.New(ioutil.Discard),
		router: router,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.Port),
			Handler: cors.Default().Handler(router),
		},
	}

	for _, option := range options {
		option(s)
	}

	s.registerRoutes()

	// We register this last so that we can use things like s.Logger inside of the `createAdminServer`
	if s.adminServer == nil {
		s.adminServer = s.createAdminServer()
	}

	return s
}

// Start starts the main web server and starts a goroutine with the admin
// server
func (s *Server) Start() error {
	go s.adminServer.ListenAndServe()
	return s.server.ListenAndServe()
}

// Shutdown calls for a graceful shutdown on the server
func (s *Server) Shutdown(ctx context.Context) error {
	s.adminServer.Shutdown(ctx)
	return s.server.Shutdown(ctx)
}

func (s *Server) createAdminServer() *http.Server {
	// Healthchecks
	h := gosundheit.New()

	// err := h.RegisterCheck(
	// 	checks.NewHostResolveCheck("api.twilio.com", 1),
	// 	gosundheit.ExecutionPeriod(60*time.Second),
	// 	gosundheit.ExecutionTimeout(2*time.Second),
	// )

	// if err != nil {
	// 	s.logger.Panic().Err(err).Msg("couldn't register healthcheck")
	// }

	mux := http.NewServeMux()
	mux.Handle("/healthz", healthhttp.HandleHealthJSON(h))

	// pprof
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	adminSrv := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.AdminPort),
		Handler: mux,
	}

	return adminSrv
}

// WithLogger sets the logger of the server
func WithLogger(logger zerolog.Logger) ServerOption {
	return func(s *Server) {
		s.logger = logger
	}
}

// WithTwilio sets the twilio client instance
func WithTwilio(twilioClient *twilio.RestClient) ServerOption {
	return func(s *Server) {
		s.twilioClient = twilioClient
	}
}

// WithDB sets the db driver
func WithDB(db *gorm.DB) ServerOption {
	return func(s *Server) {
		s.db = db
	}
}

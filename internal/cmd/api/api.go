package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/abatilo/catfacts/internal/model"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/twilio/twilio-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Cmd parses config and starts the application
func Cmd(logger zerolog.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "api",
		Short: "Runs the api web server",
		Run: func(_ *cobra.Command, _ []string) {
			cfg := &Config{
				Port:              viper.GetInt(FlagPortName),
				AdminPort:         viper.GetInt(FlagAdminPortName),
				TwilioHost:        viper.GetString(FlagTwilioHostName),
				TwilioAccountSID:  viper.GetString(FlagTwilioAccountSIDName),
				TwilioAuthToken:   viper.GetString(FlagTwilioAuthTokenName),
				TwilioPhoneNumber: viper.GetString(FlagTwilioPhoneNumberName),
				DBHost:            viper.GetString(FlagDBHost),
				DBUser:            viper.GetString(FlagDBUser),
				DBPassword:        viper.GetString(FlagDBPassword),
				DBName:            viper.GetString(FlagDBName),
				DBSSLMode:         viper.GetString(FlagDBSSLMode),
			}
			run(logger, cfg)
		}}

	cmd.PersistentFlags().Int(FlagPortName, FlagPortDefault, "The port to run the web server on")
	viper.BindPFlag(FlagPortName, cmd.PersistentFlags().Lookup(FlagPortName))

	cmd.PersistentFlags().Int(FlagAdminPortName, FlagAdminPortDefault, "The admin port to run the administrative web server on")
	viper.BindPFlag(FlagAdminPortName, cmd.PersistentFlags().Lookup(FlagAdminPortName))

	cmd.PersistentFlags().String(FlagTwilioHostName, FlagTwilioHostDefault, "Host used by Twilio webhook")
	viper.BindPFlag(FlagTwilioHostName, cmd.PersistentFlags().Lookup(FlagTwilioHostName))

	cmd.PersistentFlags().String(FlagTwilioAccountSIDName, FlagTwilioAccountSIDDefault, "Twilio account string ID")
	viper.BindPFlag(FlagTwilioAccountSIDName, cmd.PersistentFlags().Lookup(FlagTwilioAccountSIDName))

	cmd.PersistentFlags().String(FlagTwilioAuthTokenName, FlagTwilioAuthTokenDefault, "Twilio auth token")
	viper.BindPFlag(FlagTwilioAuthTokenName, cmd.PersistentFlags().Lookup(FlagTwilioAuthTokenName))

	cmd.PersistentFlags().String(FlagTwilioPhoneNumberName, FlagTwilioPhoneNumberDefault, "Twilio phone number")
	viper.BindPFlag(FlagTwilioPhoneNumberName, cmd.PersistentFlags().Lookup(FlagTwilioPhoneNumberName))

	cmd.PersistentFlags().String(FlagDBHost, FlagDBHostDefault, "DB Host")
	viper.BindPFlag(FlagDBHost, cmd.PersistentFlags().Lookup(FlagDBHost))

	cmd.PersistentFlags().String(FlagDBUser, FlagDBUserDefault, "DB User")
	viper.BindPFlag(FlagDBUser, cmd.PersistentFlags().Lookup(FlagDBUser))

	cmd.PersistentFlags().String(FlagDBPassword, FlagDBPasswordDefault, "DB Password")
	viper.BindPFlag(FlagDBPassword, cmd.PersistentFlags().Lookup(FlagDBPassword))

	cmd.PersistentFlags().String(FlagDBName, FlagDBNameDefault, "DB Name")
	viper.BindPFlag(FlagDBName, cmd.PersistentFlags().Lookup(FlagDBName))

	cmd.PersistentFlags().String(FlagDBSSLMode, FlagDBSSLModeDefault, "DB SSLMode")
	viper.BindPFlag(FlagDBSSLMode, cmd.PersistentFlags().Lookup(FlagDBSSLMode))

	return cmd
}

func run(logger zerolog.Logger, cfg *Config) {
	logger.Info().Msgf("%#v", cfg)

	// Build dependendies
	twilioClient := twilio.NewRestClient(cfg.TwilioAccountSID, cfg.TwilioAuthToken)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC", cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Panic().Err(err).Msg("Unable to connect to database")
	}

	logger.Info().Msg("Starting migrations")
	db.AutoMigrate(
		&model.Target{},
	)
	logger.Info().Msg("Finished migrations")

	// End build dependendies

	s := NewServer(cfg,
		WithLogger(logger),
		WithTwilio(twilioClient),
		WithDB(db),
	)

	// Register signal handlers for graceful shutdown
	done := make(chan struct{})
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		logger.Info().Msg("Shutting down gracefully")
		s.Shutdown(context.Background())
		close(done)
	}()

	if err := s.Start(); err != http.ErrServerClosed {
		logger.Error().Err(err).Msg("couldn't shut down gracefully")
	}
	<-done
	logger.Info().Msg("Exiting")
}

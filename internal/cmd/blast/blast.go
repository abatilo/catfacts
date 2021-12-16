package blast

import (
	"fmt"
	"time"

	"github.com/abatilo/catfacts/internal/facts"
	"github.com/abatilo/catfacts/internal/model"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/twilio/twilio-go"
	tw_api "github.com/twilio/twilio-go/rest/api/v2010"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
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
	// Twilio values
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

// Cmd parses config and starts the application
func Cmd(logger zerolog.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "blast",
		Short: "Send a Cat Fact to every active user",
		Run: func(_ *cobra.Command, _ []string) {
			cfg := &Config{
				TwilioAccountSID:  viper.GetString(FlagTwilioAccountSIDName),
				TwilioAuthToken:   viper.GetString(FlagTwilioAuthTokenName),
				TwilioPhoneNumber: viper.GetString(FlagTwilioPhoneNumberName),
				DBHost:            viper.GetString(FlagDBHost),
				DBUser:            viper.GetString(FlagDBUser),
				DBPassword:        viper.GetString(FlagDBPassword),
				DBName:            viper.GetString(FlagDBName),
				DBSSLMode:         viper.GetString(FlagDBSSLMode),
				DBSearchPath:      viper.GetString(FlagDBSearchPath),
			}
			run(logger, cfg)
		}}

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

	cmd.PersistentFlags().String(FlagDBSearchPath, FlagDBSearchPathDefault, "DB Search Path")
	viper.BindPFlag(FlagDBSearchPath, cmd.PersistentFlags().Lookup(FlagDBSearchPath))

	return cmd
}

func run(logger zerolog.Logger, cfg *Config) {
	logger.Info().Msgf("%#v", cfg)

	// Build dependendies
	twilioClient := twilio.NewRestClient(cfg.TwilioAccountSID, cfg.TwilioAuthToken)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s search_path=%s TimeZone=UTC", cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode, cfg.DBSearchPath)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Panic().Err(err).Msg("Unable to connect to database")
	}
	// End build dependendies

	var targets []model.Target
	db.Order("created_at asc").Find(&targets)

	logger.Info().Int("usersCount", len(targets)).Msg("Sending an SMS to all registered users")

	for _, target := range targets {
		timeSinceLastSMS := time.Since(target.LastSMS)

		if target.Active && 1.0 < timeSinceLastSMS.Hours() {
			randomFact := facts.GenerateFact()
			twilioClient.ApiV2010.CreateMessage(&tw_api.CreateMessageParams{
				From: &cfg.TwilioPhoneNumber,
				To:   &target.PhoneNumber,
				Body: &randomFact,
			})

			target.LastSMS = time.Now().UTC()
			db.Save(&target)
		}
	}
}

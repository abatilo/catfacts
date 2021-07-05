package main

import (
	"os"

	"github.com/abatilo/catfacts/cmd/api"
	"github.com/abatilo/catfacts/cmd/blast"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	// Define a root command to register things to
	var (
		rootCmd = &cobra.Command{
			Use:   "cf",
			Short: "Entrypoint to running the various parts of the CatFacts backend",
		}
	)

	// Prefix "CF_" when searching for environment variables
	viper.SetEnvPrefix("CF")
	// Ensure that all sub commands automatically search for env variables
	viper.AutomaticEnv()

	// Create the logger at the root of the entrypoint to guarantee it flushes
	// correctly.
	logger := zerolog.New(os.Stdout)

	rootCmd.AddCommand(api.Cmd(logger))
	rootCmd.AddCommand(blast.Cmd(logger))
	rootCmd.Execute()
}

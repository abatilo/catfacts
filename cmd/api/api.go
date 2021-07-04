package api

import (
	"github.com/abatilo/catfacts/internal/cmd/api"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

// Cmd creates the entrypoint for the api
func Cmd(logger zerolog.Logger) *cobra.Command {
	return api.Cmd(logger)
}

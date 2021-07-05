package blast

import (
	"github.com/abatilo/catfacts/internal/cmd/blast"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

// Cmd creates the entrypoint for the blast
func Cmd(logger zerolog.Logger) *cobra.Command {
	return blast.Cmd(logger)
}

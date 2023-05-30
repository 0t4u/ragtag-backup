package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get the version of rtb",
	Long:  "Get the version of the locally installed rtb binary",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msg("Version 1.0.0")
	},
}

package cmd

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	cobra.OnInitialize(initConfig)

	// rootCmd.PersistentFlags().String()
}

func initConfig() {
	// Setup ENVs
	viper.SetEnvPrefix("ragtag")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	// Config loaded from ENV second
	viper.AutomaticEnv()

	// Config loaded from flags last
	rootCmd.Flags().VisitAll(func(flag *pflag.Flag) {
		cfgName := strings.ReplaceAll(flag.Name, "-", "")

		if !flag.Changed && viper.IsSet(cfgName) {
			val := viper.Get(cfgName)
			rootCmd.Flags().Set(flag.Name, fmt.Sprintf("%v", val))
		}
	})
}

var rootCmd = &cobra.Command{
	Use:   "rtb",
	Short: "Ragtag Backup backs up data from archive.ragtag.moe",
	Long:  "A simple tool to back up data from archive.ragtag.moe by entire channel or specific video, with the option to create torrents or compress data as well.",
	Run: func(cmd *cobra.Command, args []string) {
		log.Fatal().Msg(`No operation provided. Please run with --help for more information.`)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Stack().Err(err).Msg("Something went wrong.")
	}
}

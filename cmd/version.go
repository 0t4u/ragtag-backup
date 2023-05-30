package cmd

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get the version of rtb",
	Long:  "Get the version of the locally installed rtb binary",
}

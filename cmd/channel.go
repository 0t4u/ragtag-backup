package cmd

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(channelCmd)
}

var channelCmd = &cobra.Command{
	Use: "channel",
}

package cmd

import (
	"fmt"
	"strings"
	"sync"

	"github.com/0t4u/ragtag-backup/api"
	"github.com/0t4u/ragtag-backup/backup"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(channelCmd)
}

var channelCmd = &cobra.Command{
	Use: "channel",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatal().Msgf("Expected 1 argument, got %s arguments.", fmt.Sprint(len(args)))
		}

		var id = args[0]

		result, err := api.ApiSearch(api.SearchQuery{
			ChannelId: id,
		})

		if err != nil {
			log.Fatal().Msg(err.Error())
		}

		log.Debug().Msgf("Query %v got %v results.", id, result.Hits.Total.Value)

		var exists = result.Hits.Total.Value >= 1

		if !exists {
			log.Fatal().Msgf("Channel %v does not exist.", id)
		}

		var path = rootCmd.PersistentFlags().Lookup("path").Value.String()

		var allFiles []backup.FileData

		for _, video := range result.Hits.Hits {
			var details = video.Source

			formats := strings.Split(details.FormatId, "+")

			files := getFiles(details.Files, formats)

			for _, file := range files {
				allFiles = append(allFiles, backup.FileData{
					Drive:   details.DriveBase,
					Channel: details.ChannelId,
					Video:   details.VideoId,
					File:    file,
				})
			}
		}

		var wg sync.WaitGroup
		fileChan := make(chan backup.FileData)

		wg.Add(10)

		for i := 0; i < 10; i++ {
			go backup.DownloadFile(&wg, path, fileChan)
		}

		for _, f := range allFiles {
			fileChan <- f
		}
		close(fileChan)

		wg.Wait()
	},
}

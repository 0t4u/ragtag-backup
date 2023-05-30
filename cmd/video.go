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
	rootCmd.AddCommand(videoCmd)
}

var videoCmd = &cobra.Command{
	Use: "video [VIDEOID]",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatal().Msgf("Expected 1 argument, got %s arguments.", fmt.Sprint(len(args)))
		}

		var id = args[0]

		result, err := api.ApiSearch(api.SearchQuery{
			VideoId: id,
		})

		if err != nil {
			log.Fatal().Msg(err.Error())
		}

		log.Debug().Msgf("Query %v got %v results.", id, result.Hits.Total.Value)

		var exists = result.Hits.Total.Value == 1

		if !exists {
			log.Fatal().Msgf("Video %v does not exist.", id)
		}

		var path = rootCmd.PersistentFlags().Lookup("path").Value.String()

		// entries, err := backup.ReadList(path + "/lists/local.json")
		// if err != nil {
		// 	log.Fatal().Msg("Could not read list file.")
		// }

		var details = result.Hits.Hits[0].Source

		// preExisting := false

		// for _, v := range entries {
		// 	if v.VideoId == id {
		// 		log.Warn().Msgf("Video %v already exists in list file.", id)
		// 		preExisting = true
		// 	}
		// }

		formats := strings.Split(details.FormatId, "+")

		filesTypes := getFiles(details.Files, formats)

		var files []backup.FileData

		for _, file := range filesTypes {
			files = append(files, backup.FileData{
				Drive:   details.DriveBase,
				Channel: details.ChannelId,
				Video:   details.VideoId,
				File:    file,
			})
		}

		var wg sync.WaitGroup
		fileChan := make(chan backup.FileData)

		wg.Add(10)
		for i := 0; i < 10; i++ {
			go backup.DownloadFile(&wg, path, fileChan)
		}

		for _, f := range files {
			fileChan <- f
		}
		close(fileChan)

		wg.Wait()

		// if preExisting {

		// } else {
		// 	backup.CreateListEntry(path, id, details.ChannelId)
		// }
	},
}

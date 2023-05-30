package backup

import (
	"context"
	"os"
	"path/filepath"
	"sync"

	"github.com/carlmjohnson/requests"
	"github.com/rs/zerolog/log"
	"github.com/ybbus/httpretry"
)

const ContentBaseUrl = "https://content.archive.ragtag.moe"

type FileData struct {
	Drive   string
	Channel string
	Video   string
	File    string
}

func DownloadFile(wg *sync.WaitGroup, base string, files <-chan FileData) {
	defer wg.Done()

	for file := range files {
		err := os.MkdirAll(filepath.Join(base, file.Channel, file.Video), os.ModePerm)
		if err != nil {
			log.Debug().Msgf("Could not create path at %v due to:\n", base, err)
		}

		f, err := os.Create(filepath.Join(base, file.Channel, file.Video, file.File))
		if err != nil {
			log.Debug().Msgf("Could not create file at %v due to:\n", base, err)
		}

		log.Debug().Msgf("Fetching /gd:%v/%v/%v", file.Drive, file.Video, file.File)

		cl := httpretry.NewDefaultClient()

		if err := requests.
			URL(ContentBaseUrl).
			Pathf("/gd:%v/%v/%v", file.Drive, file.Video, file.File).
			Client(cl).
			ToWriter(f).
			Fetch(context.Background()); err != nil {
			log.Error().Msgf("Could not fetch data: %v", err.Error())
		} else {
			log.Info().Msgf("Downloaded /gd:%v/%v/%v", file.Drive, file.Video, file.File)
		}

		if err = f.Close(); err != nil {
			log.Error().Msg(err.Error())
		}
	}
}

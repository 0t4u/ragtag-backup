package backup

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/rs/zerolog/log"
)

type ListEntry struct {
	VideoId   string    `json:"v"`
	ChannelId string    `json:"c"`
	Updated   time.Time `json:"u"`
}

func CreateListEntry(fileName string, video string, channel string) (err error) {
	appender, err := ArrayAppender(fileName)
	if err != nil {
		log.Debug().Msgf("Could not read file %v: %v", fileName, err.Error())
		return err
	}

	added := ListEntry{
		VideoId:   video,
		ChannelId: channel,
		Updated:   time.Now(),
	}

	if err := json.NewEncoder(appender).Encode(added); err != nil {
		log.Debug().Msgf("Could not write to file %v: %v", fileName, err.Error())
		return err
	}

	if err := appender.Close(); err != nil {
		log.Debug().Msgf("Could not close write to file %v: %v", fileName, err.Error())
		return err
	}

	return nil
}

// stolen from https://gist.github.com/rodkranz/0a8ed14fa44b5860f6668efae02b3ea5
func ReadList(fileName string) (entries []*ListEntry, err error) {
	start := time.Now()

	file, err := os.Open(fileName)
	if err != nil {
		log.Debug().Msgf("Could not open file %v: %v", fileName, err.Error())
		return nil, err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		log.Debug().Msgf("Could not stat file: %v", err.Error())
		return nil, err
	}

	reader := bufio.NewReader(file)
	decoder := json.NewDecoder(reader)

	i := 0

	var result []*ListEntry

	decoder.Token()

	for decoder.More() {
		data := &ListEntry{}
		if err := decoder.Decode(data); err != nil {
			log.Debug().Msgf("Could not decode chunk: %v", err.Error())
			return nil, err
		}

		result = append(result, data)

		i++
	}

	decoder.Token()

	elapsed := time.Since(start)

	log.Debug().Msgf("Took %v to decode %v objects from file %v with size %v", elapsed, i, fileName, fileSize(fileInfo.Size()))

	return result, nil
}

func logn(n, b float64) float64 {
	return math.Log(n) / math.Log(b)
}

func humanateBytes(s uint64, base float64, sizes []string) string {
	if s < 10 {
		return fmt.Sprintf("%dB", s)
	}
	e := math.Floor(logn(float64(s), base))
	suffix := sizes[int(e)]
	val := float64(s) / math.Pow(base, math.Floor(e))
	f := "%.0f"
	if val < 10 {
		f = "%.1f"
	}

	return fmt.Sprintf(f+"%s", val, suffix)
}

// fileSize calculates the file size and generate user-friendly string.
func fileSize(s int64) string {
	sizes := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
	return humanateBytes(uint64(s), 1024, sizes)
}

package cmd

import (
	"strings"

	"github.com/0t4u/ragtag-backup/api"
)

func getFiles(files []api.VideoFiles, formats []string) []string {
	s := []string{}

	for _, v := range files {
		for _, f := range formats {
			if strings.Contains(v.Name, f) {
				s = append(s, v.Name)
			}
			if strings.Contains(v.Name, ".json") {
				s = append(s, v.Name)
			}
			if strings.Contains(v.Name, ".webp") {
				s = append(s, v.Name)
			}
		}
	}

	return s
}

package api_test

import (
	"testing"

	"github.com/0t4u/ragtag-backup/api"
)

func TestSearch(t *testing.T) {
	json, err := api.ApiSearch(api.SearchQuery{
		VideoId: "D9bCf2RLRVY",
	})
	if err != nil {
		t.Error(err)
	}

	if json.Hits.Total.Value == 0 {
		t.Error("API returned 0 results")
	}

	if json.Hits.Hits[0].Id != "D9bCf2RLRVY" {
		t.Error("VideoId does not match")
	}
}

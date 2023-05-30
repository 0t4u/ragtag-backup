package backup_test

import (
	"testing"

	"github.com/0t4u/ragtag-backup/backup"
)

func TestReadList(t *testing.T) {
	data, err := backup.ReadList("./sample.json")

	if err != nil {
		t.Error(err)
	}

	if data[len(data)-1].ChannelId != "4" {
		t.Errorf("Result %v does not match expected result \"3\"", data[len(data)-1].ChannelId)
	}
}

func TestWriteEntry(t *testing.T) {
	err := backup.CreateListEntry("./sample.json", "4", "4")

	if err != nil {
		t.Error(err)
	}
}

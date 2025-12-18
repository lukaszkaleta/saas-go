package filestore

import (
	"log/slog"
	"os"
	"testing"
)

func TestFileRecordModel(t *testing.T) {
	file, err := os.OpenFile("./record-api.go", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		t.Error(err)
	}
	model := FileRecordModel(*file)
	slog.Info(model.Url)
}

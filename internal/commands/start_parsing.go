package commands

import (
	"context"
	"log"
	"os"

	"github.com/HugoDrl/zebra/internal/parser"
	"github.com/HugoDrl/zebra/internal/reader"
	"github.com/HugoDrl/zebra/internal/store"
)

func StartFileParsing(ctx context.Context, d store.LogsStore, filename string, json bool) error {
	root, err := os.OpenRoot(".")
	if err != nil {
		return nil
	}

	filectx, cancel := context.WithCancel(ctx)
	logsChan, errsChan := reader.ExtractLogsFromFileName(reader.ExtractLinesFromFileInput{
		Ctx:           filectx,
		Root:          root,
		Filename:      filename,
		ParseFunction: parser.GetParseFunction(parser.ParseSettings{Json: json}),
	})

	go func() {
		if _, err := d.InsertLogs(logsChan, errsChan); err != nil {
			cancel()
		}
	}()
	go func() {
		for err := range errsChan {
			log.Println(err)
		}
	}()

	return nil
}

package commands

import (
	"os"

	"github.com/HugoDrl/zebra/internal/data"
	"github.com/HugoDrl/zebra/internal/parser"
	"github.com/HugoDrl/zebra/internal/reader"
)

func StartFileParsing(d *data.DataLayer, filename string, json bool) error {
	root, err := os.OpenRoot(".")
	if err != nil {
		return err
	}

	logsChan, errsChan := reader.ExtractLogsFromFileName(reader.ExtractLinesFromFileInput{
		Root:          root,
		Filename:      filename,
		ParseFunction: parser.GetParseFunction(parser.ParseSettings{Json: json}),
	})

	go d.IngestLogChan(logsChan)
	go d.IngestErrChan(errsChan)

	return nil
}

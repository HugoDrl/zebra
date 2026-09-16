package server

import (
	"net/http"
	"os"

	"github.com/HugoDrl/zebra/internal/parser"
	"github.com/HugoDrl/zebra/internal/reader"
)

func (d *DataLayer) startFileParsing(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	if filename == "" {
		w.WriteHeader(400)
		return
	}

	root, err := os.OpenRoot(".")
	if err != nil {
		w.WriteHeader(500)
		return
	}
	logsChan, errsChan := reader.ExtractLogsFromFileName(reader.ExtractLinesFromFileInput{
		Root:          root,
		Filename:      filename,
		ParseFunction: parser.ParseJSONFormatLine,
	})

	go d.ingestLogChan(logsChan)
	go d.ingestErrChan(errsChan)
	w.WriteHeader(200)
}

func (d *DataLayer) AttachCommandsHandler(handler *http.ServeMux) {
	handler.HandleFunc("GET /parse", d.startFileParsing)
}

package server

import (
	"net/http"
	"strings"

	"github.com/HugoDrl/zebra/internal/commands"
)

func (s *HttpServer) startFileParsing(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	if filename == "" {
		w.WriteHeader(400)
		return
	}
	jsonQuery := strings.ToLower(r.URL.Query().Get("json"))
	json := jsonQuery == "true" || jsonQuery == "1"

	err := commands.StartFileParsing(s.dataLayer, filename, json)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
}

func (s *HttpServer) AttachCommandsHandler(handler *http.ServeMux) {
	handler.HandleFunc("GET /parse", s.startFileParsing)
}

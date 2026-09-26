package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/HugoDrl/zebra/internal/analyser"
	"github.com/HugoDrl/zebra/internal/filter"
)

func (s *HttpServer) retrieveSlowestLogs(w http.ResponseWriter, r *http.Request) {
	numberOfSlowestLogs, err := strconv.Atoi(r.URL.Query().Get("number_of_logs"))
	if err != nil {
		w.WriteHeader(422)
		return
	}
	slowestLogs := analyser.RetrieveSlowestLogs(s.dataLayer.Logs, numberOfSlowestLogs)
	payload, err := json.Marshal(slowestLogs)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	if _, err := w.Write(payload); err != nil {
		w.WriteHeader(500)
		return
	}
}

func (s *HttpServer) retrieveLogs(w http.ResponseWriter, r *http.Request) {
	filters, err := filter.ProcessRequestToFilter(*r)
	if err != nil {
		w.WriteHeader(422)
		return
	}

	payload, err := json.Marshal(filter.FilterLogs(s.dataLayer.Logs, filters))
	if err != nil {
		w.WriteHeader(500)
		return
	}

	if _, err := w.Write(payload); err != nil {
		w.WriteHeader(500)
		return
	}
}

func (s *HttpServer) AttachLogsHandler(mux *http.ServeMux) {
	mux.HandleFunc("GET /logs", s.retrieveLogs)
	mux.HandleFunc("GET /slowest-logs", s.retrieveSlowestLogs)
}

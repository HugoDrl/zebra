package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/HugoDrl/zebra/internal/analyser"
)

func (d *DataLayer) retrieveSlowestLogs(w http.ResponseWriter, r *http.Request) {
	numberOfSlowestLogs, err := strconv.Atoi(r.URL.Query().Get("number_of_logs"))
	if err != nil {
		w.WriteHeader(422)
		return
	}
	slowestLogs := analyser.RetrieveSlowestLogs(d.Logs, numberOfSlowestLogs)
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

func (d *DataLayer) retrieveLogs(w http.ResponseWriter, r *http.Request) {
	payload, err := json.Marshal(d.Logs)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	if _, err := w.Write(payload); err != nil {
		w.WriteHeader(500)
		return
	}
}

func (d *DataLayer) AttachLogsHandler(mux *http.ServeMux) {
	mux.HandleFunc("GET /logs", d.retrieveLogs)
	mux.HandleFunc("GET /slowest-logs", d.retrieveSlowestLogs)
}

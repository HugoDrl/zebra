package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/HugoDrl/zebra/internal/analyser"
)

func (d *DataLayer) retrieveLogs(w http.ResponseWriter, r *http.Request) {
	payload, err := json.Marshal(d.Logs)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	if _, err := w.Write(payload); err != nil {
		w.WriteHeader(500)
	}
}

func (d *DataLayer) calculateMetrics(w http.ResponseWriter, r *http.Request) {
	metrics := analyser.AnalyseLogs(d.Logs, d.Errs)

	payload, err := json.Marshal(metrics)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	if _, err := w.Write(payload); err != nil {
		w.WriteHeader(500)
	}
}

func (d *DataLayer) retrieveSlowestLogs(w http.ResponseWriter, r *http.Request) {
	slowestLogsNumber, err := strconv.Atoi(r.URL.Query().Get("number_of_logs"))
	if err != nil {
		w.WriteHeader(422)
		if _, err := w.Write([]byte("number_of_logs is mandatory and should be a number")); err != nil {
			w.WriteHeader(500)
		}
		return
	}

	slowestLogs := analyser.RetrieveSlowestLogs(d.Logs, slowestLogsNumber)
	payload, err := json.Marshal(slowestLogs)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	if _, err := w.Write(payload); err != nil {
		w.WriteHeader(500)
	}
}

func (d *DataLayer) AttachMetricsHandler(handler *http.ServeMux) {
	handler.HandleFunc("GET /logs", d.retrieveLogs)
	handler.HandleFunc("GET /metrics", d.calculateMetrics)
	handler.HandleFunc("GET /slowest_logs", d.retrieveSlowestLogs)
}

package server

import (
	"encoding/json"
	"net/http"

	"github.com/HugoDrl/zebra/internal/analyser"
)

func (s *HttpServer) calculateMetrics(w http.ResponseWriter, r *http.Request) {
	metrics := analyser.AnalyseLogs(s.dataLayer.Logs, s.dataLayer.Errs)

	payload, err := json.Marshal(metrics)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	if _, err := w.Write(payload); err != nil {
		w.WriteHeader(500)
	}
}

func (s *HttpServer) AttachMetricsHandler(handler *http.ServeMux) {
	handler.HandleFunc("GET /metrics", s.calculateMetrics)
}

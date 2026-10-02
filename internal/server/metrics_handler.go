package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/HugoDrl/zebra/internal/analyser"
	"github.com/HugoDrl/zebra/internal/filter"
)

func (s *HttpServer) calculateMetrics(w http.ResponseWriter, r *http.Request) {
	filters, err := filter.ProcessRequestToFilter(*r)
	if err != nil {
		w.WriteHeader(422)
		return
	}

	logs, errs, err := s.LogStore.RetrieveLogs(filters)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}
	metrics := analyser.AnalyseLogs(logs, errs)

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

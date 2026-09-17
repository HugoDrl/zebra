package server

import (
	"net/http"
	"time"

	"github.com/HugoDrl/zebra/internal/data"
)

type HttpServer struct {
	server    *http.Server
	dataLayer *data.DataLayer
}

func NewServer(d *data.DataLayer) *HttpServer {
	handler := http.NewServeMux()

	s := &HttpServer{
		server:    &http.Server{Addr: ":8000", ReadHeaderTimeout: 100 * time.Millisecond},
		dataLayer: d,
	}
	s.AttachCommandsHandler(handler)
	s.AttachLogsHandler(handler)
	s.AttachMetricsHandler(handler)
	s.server.Handler = handler
	return s
}

func (s *HttpServer) StartServer() error {
	return s.server.ListenAndServe()
}

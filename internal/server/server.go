package server

import (
	"net/http"
	"time"

	"github.com/HugoDrl/zebra/internal/parser"
)

type DataLayer struct {
	Logs []*parser.Log
	Errs []error
}

func (d *DataLayer) ingestLogChan(c <-chan *parser.Log) {
	for log := range c {
		d.Logs = append(d.Logs, log)
	}
}

func (d *DataLayer) ingestErrChan(c <-chan error) {
	for err := range c {
		d.Errs = append(d.Errs, err)
	}
}

type HttpServer struct {
	server    *http.Server
	dataLayer *DataLayer
}

func NewServer(data *DataLayer) *HttpServer {
	handler := http.NewServeMux()
	data.AttachCommandsHandler(handler)
	data.AttachLogsHandler(handler)
	data.AttachMetricsHandler(handler)

	s := &HttpServer{
		server:    &http.Server{Addr: ":8000", ReadHeaderTimeout: 100 * time.Millisecond},
		dataLayer: data,
	}
	s.server.Handler = handler
	return s
}

func (s *HttpServer) StartServer() error {
	return s.server.ListenAndServe()
}

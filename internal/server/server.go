package server

import (
	"context"
	"net/http"
	"time"

	"github.com/HugoDrl/zebra/internal/store"
)

type HttpServer struct {
	server   *http.Server
	ctx      context.Context
	cancel   func()
	LogStore store.LogsStore
	Errs     []error
}

func NewServer(logStore store.LogsStore) *HttpServer {
	handler := http.NewServeMux()
	ctx, cancel := context.WithCancel(context.Background())

	s := &HttpServer{
		server:   &http.Server{Addr: ":8000", ReadHeaderTimeout: 100 * time.Millisecond},
		ctx:      ctx,
		cancel:   cancel,
		LogStore: logStore,
		Errs:     make([]error, 0),
	}
	s.AttachCommandsHandler(handler)
	s.AttachLogsHandler(handler)
	s.AttachMetricsHandler(handler)
	s.server.Handler = handler
	return s
}

func (s *HttpServer) StartServer() {
	go s.server.ListenAndServe()
}

func (s *HttpServer) Cancel() {
	s.cancel()
}

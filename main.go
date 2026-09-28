package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/HugoDrl/zebra/internal/data"
	"github.com/HugoDrl/zebra/internal/server"
)

func main() {
	ctx := make(chan os.Signal, 1)
	signal.Notify(ctx, syscall.SIGINT, syscall.SIGTERM)

	dataLayer, cancel := data.NewDataLayer()

	s := server.NewServer(
		dataLayer,
	)
	s.StartServer()

	<-ctx
	cancel()
}

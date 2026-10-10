package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/HugoDrl/zebra/internal/server"
	"github.com/HugoDrl/zebra/internal/store"
)

func main() {
	ctx := make(chan os.Signal, 1)
	signal.Notify(ctx, syscall.SIGINT, syscall.SIGTERM)

	store, err := store.NewSQLiteStore("./temp.db")
	if err != nil {
		log.Fatal(err)
	}

	s := server.NewServer(
		store,
	)
	s.StartServer()

	log.Println("succesfully started server on port 8000")

	<-ctx
	s.Cancel()
	log.Println("server stopped")
}

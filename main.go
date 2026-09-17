package main

import (
	"fmt"
	"os"

	"github.com/HugoDrl/zebra/internal/data"
	"github.com/HugoDrl/zebra/internal/server"
)

func main() {
	if err := server.NewServer(
		data.NewDataLayer(),
	).StartServer(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

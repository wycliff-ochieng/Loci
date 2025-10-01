package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/wycliff-ochieng/api"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	fmt.Println("Spinning up the LOCI APPLICATION")
	server := api.NewServer(logger, ":3000")
	server.Run()
}

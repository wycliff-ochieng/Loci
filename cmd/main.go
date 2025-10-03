package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/wycliff-ochieng/api"
	"github.com/wycliff-ochieng/internal/config"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	conf, err := config.Load()
	if err != nil {
		log.Fatal("failed to load env")
	}
	fmt.Println("Spinning up the LOCI APPLICATION")
	server := api.NewServer(logger, ":3000", conf)
	server.Run()
}

package api

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/wycliff-ochieng/internal/config"
	"github.com/wycliff-ochieng/internal/store"
)

type Server struct {
	l    *slog.Logger
	addr string
	cfg  *config.Config
}

func NewServer(l *slog.Logger, addr string) *Server {
	return &Server{
		l:    l,
		addr: addr,
	}
}

func LociTest(w http.ResponseWriter, r *http.Request) {
	fmt.Println("TEsting the server")
}

func (s *Server) Run() {
	fmt.Println("starting up run method")

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)

	defer cancel()

	//handler

	//service

	db, err := store.NewPostgis(ctx, s.cfg)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	router := mux.NewRouter()

	test := router.Methods("POST").Subrouter()
	test.HandleFunc("/", LociTest)

	if err := http.ListenAndServe(s.addr, router); err != nil {
		fmt.Errorf("rror listening to server: %s", err)
	}

}

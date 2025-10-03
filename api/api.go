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
	"github.com/wycliff-ochieng/internal/service"
	"github.com/wycliff-ochieng/internal/store"
	handlers "github.com/wycliff-ochieng/internal/transport/http"
	"github.com/wycliff-ochieng/sqlc"
)

type Server struct {
	l    *slog.Logger
	addr string
	cfg  *config.Config
}

func NewServer(l *slog.Logger, addr string, cfg *config.Config) *Server {
	return &Server{
		l:    l,
		addr: addr,
		cfg:  cfg,
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

	db, err := store.NewPostgis(ctx, s.cfg)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	queries := sqlc.New(db)

	us := service.NewUserService(db, *queries)

	//handler
	uh := handlers.NewUserHandler(logger, us)

	router := mux.NewRouter()

	test := router.Methods("POST").Subrouter()
	test.HandleFunc("/register", uh.Register)

	if err := http.ListenAndServe(s.addr, router); err != nil {
		fmt.Errorf("rror listening to server: %s", err)
	}

}

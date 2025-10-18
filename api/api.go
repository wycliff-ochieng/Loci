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

	register := router.Methods("POST").Subrouter()
	register.HandleFunc("/register", uh.Register)

	login := router.Methods("POST").Subrouter()
	login.HandleFunc("/login", uh.Login)

	getLoci := router.Methods("GET").Subrouter()
	getLoci.HandleFunc("/api/get/loci/{location}", uh.GetLociInGeoFencedLocation)

	if err := http.ListenAndServe(s.addr, router); err != nil {
		log.Fatalf("rror listening to server: %s", err)
	}

}

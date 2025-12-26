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
	"github.com/redis/go-redis/v9"
	"github.com/wycliff-ochieng/internal/config"
	"github.com/wycliff-ochieng/internal/limitter"
	"github.com/wycliff-ochieng/internal/service"
	"github.com/wycliff-ochieng/pkg/middleware"

	//"github.com/wycliff-ochieng/internal/socket"
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

	rclt := redis.NewClient(&redis.Options{
		Addr:     s.cfg.REDIS_ADDR,
		Password: s.cfg.REDIS_PASSWORD,
		DB:       0,
	})

	wdw := 60 * time.Second

	lim := int16(20)

	rtl := limitter.NewRedisLimitter(rclt, wdw, lim)

	queries := sqlc.New(db)

	authMiddleware := middleware.AuthenticationMiddleware(s.cfg.JWTsecret)

	us := service.NewUserService(db, *queries, rtl)

	//handler
	uh := handlers.NewUserHandler(logger, us)

	router := mux.NewRouter()

	register := router.Methods("POST").Subrouter()
	register.HandleFunc("/register", uh.Register)

	login := router.Methods("POST").Subrouter()
	login.HandleFunc("/login", uh.Login)

	getLoci := router.Methods("GET").Subrouter()
	getLoci.HandleFunc("/api/get/loci/", uh.GetLociInGeoFencedLocation)
	getLoci.Use(authMiddleware)

	postLoci := router.Methods("POST").Subrouter()
	postLoci.HandleFunc("/api/post/loci", uh.CreateLoci)
	postLoci.Use(authMiddleware)

	//http.HandleFunc("/ws",socket.ServerWS)

	if err := http.ListenAndServe(s.addr, router); err != nil {
		log.Fatalf("rror listening to server: %s", err)
	}

}

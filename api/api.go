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
	"github.com/wycliff-ochieng/internal/socket"
	"github.com/wycliff-ochieng/pkg/middleware"

	//"github.com/wycliff-ochieng/internal/socket"
	corshandlers "github.com/gorilla/handlers"
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

	//initialize hub
	hub := socket.NewHub()

	go hub.Run()

	us := service.NewUserService(db, *queries, rtl, hub)

	//handler
	uh := handlers.NewUserHandler(logger, us)

	router := mux.NewRouter()
	//srouter.Use(corsMiddleware(s.cfg.CORSOrigin))

	register := router.Methods("POST").Subrouter()
	register.HandleFunc("/register", uh.Register)

	login := router.Methods("POST").Subrouter()
	login.HandleFunc("/login", uh.Login)

	getLoci := router.Methods("GET").Subrouter()
	getLoci.HandleFunc("/api/get/loci/", uh.GetLociInGeoFencedLocation)
	// public read for map

	postLoci := router.Methods("POST").Subrouter()
	postLoci.HandleFunc("/api/post/loci", uh.CreateLoci)
	postLoci.Use(authMiddleware)

	viewLoci := router.Methods("POST").Subrouter()
	viewLoci.HandleFunc("/api/loci/{id}/view", uh.ViewLociHandler)
	viewLoci.Use(authMiddleware)

	replyLoci := router.Methods("POST").Subrouter()
	replyLoci.HandleFunc("/api/loci/{id}/reply", uh.ReplyToLociHandler)
	replyLoci.Use(authMiddleware)

	getReplies := router.Methods("GET").Subrouter()
	getReplies.HandleFunc("/api/loci/{id}/replies", uh.GetRepliesHandler)
	// public read for threads

	wsRouter := router.Methods("GET").Subrouter()
	wsRouter.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		socket.ServerWS(hub, w, r)
	})
	wsRouter.Use(authMiddleware)

	//CORS configurtion
	origins := s.cfg.CORSAllowedOrigins

	allowedMethods := corshandlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	allowedHeaders := corshandlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	allowCredentials := corshandlers.AllowCredentials()
	allowedOrigins := corshandlers.AllowedOrigins(origins)

	cm := corshandlers.CORS(allowedOrigins, allowCredentials, allowedMethods, allowedHeaders)(router)

	if err := http.ListenAndServe(s.addr, cm); err != nil {
		log.Fatalf("rror listening to server: %s", err)
	}

}

// corsMiddleware allows the frontend (Next dev on 3001 by default) to reach the Go API on 3000.
func corsMiddleware(origin string) mux.MiddlewareFunc {
	allowed := map[string]struct{}{
		origin:                  {},
		"http://localhost:3001": {},
		"http://127.0.0.1:3001": {},
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqOrigin := r.Header.Get("Origin")
			if _, ok := allowed[reqOrigin]; ok {
				w.Header().Set("Access-Control-Allow-Origin", reqOrigin)
				w.Header().Set("Vary", "Origin")
			}
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Origin, Authorization, Content-Type, Accept")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

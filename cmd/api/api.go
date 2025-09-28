package api

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	l    *slog.Logger
	addr string
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

	router := mux.NewRouter()

	test := router.Methods("POST").Subrouter()
	test.HandleFunc("/", LociTest)

	if err := http.ListenAndServe(s.addr, router); err != nil {
		fmt.Errorf("rror listening to server: %s", err)
	}

}

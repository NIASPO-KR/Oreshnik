package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"oreshnik/config"
	"oreshnik/internal/datacontroller"
)

type Server struct {
	cfg *config.Config

	dc *datacontroller.DataController

	router *chi.Mux
	server *http.Server
}

func New(cfg *config.Config) (*Server, error) {
	s := &Server{
		cfg: cfg,
	}

	if err := s.init(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Server) init() error {
	s.dc = datacontroller.New(
		fmt.Sprintf("%s:%s", s.cfg.Microservices.Static.Addr, s.cfg.Microservices.Static.Port),
		fmt.Sprintf("%s:%s", s.cfg.Microservices.Users.Addr, s.cfg.Microservices.Users.Port),
		newHTTPClient(),
	)

	s.initRouter()
	s.initHTTPServer()

	return nil
}

func (s *Server) initHTTPServer() {
	s.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%s", s.cfg.Server.Addr, s.cfg.Server.Port),
		Handler:      s.router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

func (s *Server) Run() {
	log.Println("Server started")

	if err := s.server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

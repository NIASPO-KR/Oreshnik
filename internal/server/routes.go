package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (s *Server) initRouter() {
	s.router = chi.NewRouter()

	s.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Content-Length"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	s.router.Route("/api", func(r chi.Router) {
		r.Route("/static", s.registerStaticRoutes)
		r.Route("/users", s.registerUsersRoutes)
	})
}

func (s *Server) registerStaticRoutes(r chi.Router) {
	r.Get("/items", s.dc.GetItems)
	r.Get("/pickupPoints", s.dc.GetPickupPoints)
	r.Get("/payments", s.dc.GetPayments)
}

func (s *Server) registerUsersRoutes(r chi.Router) {
	r.Get("/cart", s.dc.GetCart)
}

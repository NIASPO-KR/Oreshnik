package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (s *Server) initRouter() {
	s.router = chi.NewRouter()

	s.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "OPTIONS"},
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
	r.Route("/cart", func(r chi.Router) {
		r.Patch("/", s.dc.UpdateCart)
		r.Get("/", s.dc.GetCart)
	})
	r.Route("/orders", func(r chi.Router) {
		r.Post("/", s.dc.CreateOrder)
		r.Patch("/", s.dc.UpdateOrderStatus)
		r.Get("/", s.dc.GetOrders)
	})
}

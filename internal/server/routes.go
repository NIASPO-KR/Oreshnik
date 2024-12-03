package server

import "github.com/go-chi/chi/v5"

func (s *Server) initRouter() {
	s.router = chi.NewRouter()

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

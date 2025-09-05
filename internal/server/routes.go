package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (s *Server) registerRoutes(r *chi.Mux) {
	// Health check
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	// SSE Notification endpoint
	r.Get("/notifications", s.NotificationsHandler)
	r.Post("/send", s.SendNotificationHandler)

	// Metrics endpoint
	r.Handle("/metrics", promhttp.Handler())
}

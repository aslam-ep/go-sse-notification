package server

import (
	"context"
	"log"
	"net/http"

	"github.com/aslam-ep/go-sse-notification/config"
	"github.com/aslam-ep/go-sse-notification/internal/notifications"
	"github.com/aslam-ep/go-sse-notification/internal/redis"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	Config  *config.Config
	Manager *notifications.Manager
	Redis   *redis.Client
	History *redis.History
}

// NewServer builds the server with dependencies
func NewServer(cfg *config.Config) *Server {
	redisClient := redis.NewClient(cfg.Redis.URL, "", 1)

	return &Server{
		Config:  cfg,
		Manager: notifications.NewManager(),
		Redis:   redisClient,
		History: redis.NewHistory(redisClient),
	}
}

// Start boots background workers + HTTP server
func (s *Server) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start background subscriber worker
	go s.startSubscriber(ctx)

	// Setup router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	s.registerRoutes(r)

	address := s.Config.Server.Address
	log.Printf("HTTP server listening on %s\n", address)
	return http.ListenAndServe(address, r)
}

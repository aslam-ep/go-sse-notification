package server

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

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
	redisClient := redis.NewClient(cfg.Redis.URL, cfg.Redis.Password, cfg.Redis.DB)

	return &Server{
		Config:  cfg,
		Manager: notifications.NewManager(),
		Redis:   redisClient,
		History: redis.NewHistory(redisClient),
	}
}

// Start boots background workers + HTTP server
func (s *Server) Start() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Start background subscriber worker
	go s.startSubscriber(ctx)

	// Setup router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	s.registerRoutes(r)

	srv := &http.Server{
		Addr:    s.Config.Server.Address,
		Handler: r,
	}

	go func() {
		<-ctx.Done()
		log.Println("Shutting down HTTP server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("Gracefull shutdown failed: %v", err)
		}
	}()

	log.Printf("HTTP server listening on %s\n", s.Config.Server.Address)
	return srv.ListenAndServe()
}

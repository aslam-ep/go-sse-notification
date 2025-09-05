package server

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aslam-ep/go-sse-notification/internal/notifications"
)

// startSubscriber listens to Redis pubsub and delivers to connected clients
func (s *Server) startSubscriber(ctx context.Context) {
	sub := s.Redis.Subscribe(ctx, "notifications")
	ch := sub.Channel()

	go func() {
		for {
			select {
			case msg := <-ch:
				var n notifications.Notification
				if err := json.Unmarshal([]byte(msg.Payload), &n); err != nil {
					log.Printf("Failed to parse pubsub: %v\n", err)
					continue
				}
				s.Manager.Send(n) // Fan out to active SSE clients
			case <-ctx.Done():
				_ = sub.Close()
				return
			}
		}
	}()
}

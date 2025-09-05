package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/aslam-ep/go-sse-notification/internal/notifications"
	"github.com/aslam-ep/go-sse-notification/internal/redis"
)

func (s *Server) NotificationsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userID")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Parse lastEventID (from header or query)
	var lastSeenID int64
	if id := r.Header.Get("Last-Event-ID"); id != "" {
		lastSeenID, _ = strconv.ParseInt(id, 10, 64)
	} else if id := r.URL.Query().Get("lastEventID"); id != "" {
		lastSeenID, _ = strconv.ParseInt(id, 10, 64)
	}

	// SSE Headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE is not supported", http.StatusInternalServerError)
		return
	}

	// reply missed messages
	if lastSeenID > 0 {
		missed, err := s.History.FetchSince(r.Context(), userID, lastSeenID)
		if err == nil {
			for _, m := range missed {
				fmt.Fprintf(w, "id: %d\ndata: %s\n\n", m.ID, m.Message)
				flusher.Flush()
			}
		} else {
			log.Printf("Failed to load the history for user %s: %v\n", userID, err)
		}
	}

	// Channel for message
	ch := make(chan notifications.Notification, 1) // Buffered so we don't block
	s.Manager.AddClient(userID, ch)

	// Cleanup when client disconnects
	defer func() {
		s.Manager.RemoveClient(userID, ch)
		close(ch)
		log.Printf("Connection closed for user %s\n", userID)
	}()

	notify := r.Context().Done()
	for {
		select {
		case n := <-ch:
			fmt.Fprintf(w, "id: %d\ndata: %s\n\n", n.ID, n.Message)
			flusher.Flush()

		case <-notify:
			// Context cancled, client disconnected
			return
		}
	}
}

// Request payload for /send
type sendRequest struct {
	UserID  string `json:"userID"`
	Message string `json:"message"`
}

// POST /send
func (s *Server) SendNotificationHandler(w http.ResponseWriter, r *http.Request) {
	var req sendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Inavlid request body", http.StatusBadRequest)
		return
	}

	if req.UserID == "" || req.Message == "" {
		http.Error(w, "Missing userID or message", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	id, err := s.History.NextID(ctx)
	if err != nil {
		http.Error(w, "Failed to generate notifiation ID", http.StatusInternalServerError)
		log.Printf("Failed to generate Notification ID: %v\n", err)
		return
	}

	n := notifications.Notification{
		ID:      id,
		UserID:  req.UserID,
		Message: req.Message,
	}

	if err := s.History.Store(ctx, n); err != nil {
		http.Error(w, "Failed to persist notifcations", http.StatusInternalServerError)
		log.Printf("Redis.History.Store: %s\n", err)
		return
	}

	payload, _ := json.Marshal(n)
	if err := s.Redis.Publish(ctx, redis.NotificationTopic, string(payload)); err != nil {
		http.Error(w, "failed to push notification", http.StatusInternalServerError)
		log.Printf("Redis Push: %s\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Notification sent"))
}

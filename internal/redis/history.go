package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aslam-ep/go-sse-notification/internal/notifications"
)

type History struct {
	rdb *Client
}

func NewHistory(c *Client) *History {
	return &History{rdb: c}
}

func (h *History) Store(ctx context.Context, n notifications.Notification) error {
	key := fmt.Sprintf("%s:%s:history", NotificationTopic, n.UserID)

	message, err := json.Marshal(n)
	if err != nil {
		return err
	}

	pipe := h.rdb.RDB.TxPipeline()

	pipe.LPush(ctx, key, message)
	pipe.LTrim(ctx, key, 0, 99)
	pipe.Expire(ctx, key, 24*time.Hour)

	_, err = pipe.Exec(ctx)

	return err
}

func (h *History) FetchSince(ctx context.Context, userID string, lastSeenID int64) ([]notifications.Notification, error) {
	key := fmt.Sprintf("%s:%s:history", NotificationTopic, userID)

	raw, err := h.rdb.RDB.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	var result []notifications.Notification
	for _, r := range raw {
		var n notifications.Notification
		if err = json.Unmarshal([]byte(r), &n); err != nil {
			continue
		}
		if n.ID > lastSeenID {
			result = append(result, n)
		}
	}

	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return result, nil
}

func (h *History) NextID(ctx context.Context) (int64, error) {
	return h.rdb.RDB.Incr(ctx, fmt.Sprintf("%s:seq", NotificationTopic)).Result()
}

package cache_key

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	CacheTtl5Mins = 5 * time.Minute
)

func CacheKeyUserByID(id uuid.UUID) string {
	return fmt.Sprintf("[STORYWEB]:user:user_id:%s", id.String())
}

func CacheKeyUserByUsername(username string) string {
	return fmt.Sprintf("[STORYWEB]:user:username:%s", username)
}

func CacheKeyUserByEmail(email string) string {
	return fmt.Sprintf("[STORYWEB]:user:email:%s", email)
}

package middleware

import (
	"net/http"
	"sync"
	"time"

	apperrors "github.com/Kal-el21/backend/internal/shared/errors"
	"github.com/Kal-el21/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
)

type rateLimiter struct {
	mu          sync.Mutex
	attempts    map[string][]time.Time
	maxAttempts int
	window      time.Duration
}

// NewRateLimiter generik, dipakai untuk login, refresh-token, dan create-transaction.
func NewRateLimiter(maxAttempts int, window time.Duration) *rateLimiter {
	rl := &rateLimiter{
		attempts:    make(map[string][]time.Time),
		maxAttempts: maxAttempts,
		window:      window,
	}

	go func() {
		ticker := time.NewTicker(window)
		for range ticker.C {
			rl.cleanup()
		}
	}()

	return rl
}

func (rl *rateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	for key, times := range rl.attempts {
		valid := []time.Time{}
		for _, t := range times {
			if now.Sub(t) < rl.window {
				valid = append(valid, t)
			}
		}
		if len(valid) == 0 {
			delete(rl.attempts, key)
		} else {
			rl.attempts[key] = valid
		}
	}
}

// LimitByIP membatasi percobaan per IP address, generik untuk endpoint apa pun.
func (rl *rateLimiter) LimitByIP() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP()
		rl.check(c, key)
	}
}

// LimitByUser membatasi percobaan per user_id (dipanggil setelah AuthMiddleware,
// sehingga user_id sudah ada di context).
func (rl *rateLimiter) LimitByUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint64("user_id")
		key := "user:" + string(rune(userID)) // identifier unik per user
		rl.check(c, key)
	}
}

func (rl *rateLimiter) check(c *gin.Context, key string) {
	rl.mu.Lock()
	now := time.Now()
	attempts := rl.attempts[key]

	valid := []time.Time{}
	for _, t := range attempts {
		if now.Sub(t) < rl.window {
			valid = append(valid, t)
		}
	}

	if len(valid) >= rl.maxAttempts {
		rl.mu.Unlock()
		response.Error(c, apperrors.New(apperrors.ErrRateLimited, "too many requests, please try again later"))
		c.AbortWithStatus(http.StatusTooManyRequests)
		return
	}

	valid = append(valid, now)
	rl.attempts[key] = valid
	rl.mu.Unlock()

	c.Next()
}

// Backward-compatible alias untuk kode Phase 1 yang sudah memanggil
// NewLoginRateLimiter dan LoginRateLimit().
func NewLoginRateLimiter(maxAttempts int, window time.Duration) *rateLimiter {
	return NewRateLimiter(maxAttempts, window)
}

func (rl *rateLimiter) LoginRateLimit() gin.HandlerFunc {
	return rl.LimitByIP()
}

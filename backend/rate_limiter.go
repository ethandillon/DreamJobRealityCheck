package main

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// RateLimiter provides simple in-memory IP based rate limiting.
// It is NOT distributed and should only be used for small scale deployments.
type RateLimiter struct {
	limit       int           // max requests per window
	window      time.Duration // e.g. 1 minute
	mu          sync.Mutex
	visitors    map[string]*visitor
	lastCleanup time.Time
}

type visitor struct {
	count       int
	windowStart time.Time
	lastSeen    time.Time
}

// NewRateLimiter creates a new rate limiter.
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		limit:    limit,
		window:   window,
		visitors: make(map[string]*visitor),
	}
}

// Middleware returns an http middleware enforcing the rate limit.
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := rl.extractIP(r)
		if rl.exceeded(ip) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Retry-After", "60") // rough seconds remainder
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"error":"rate limit exceeded"}`))
			return
		}
		next.ServeHTTP(w, r)
	})
}

// exceeded increments the counter for IP and returns true if over limit.
func (rl *RateLimiter) exceeded(ip string) bool {
	now := time.Now()
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.lastCleanup.IsZero() || now.Sub(rl.lastCleanup) > 5*rl.window {
		for k, v := range rl.visitors {
			if now.Sub(v.lastSeen) > 10*rl.window { // stale
				delete(rl.visitors, k)
			}
		}
		rl.lastCleanup = now
	}

	v, ok := rl.visitors[ip]
	if !ok {
		v = &visitor{count: 0, windowStart: now, lastSeen: now}
		rl.visitors[ip] = v
	}
	v.lastSeen = now

	if now.Sub(v.windowStart) >= rl.window {
		v.count = 0
		v.windowStart = now
	}
	v.count++
	return v.count > rl.limit
}

// extractIP attempts to determine the real client IP accounting for proxies.
func (rl *RateLimiter) extractIP(r *http.Request) string {
	hdrs := []string{"X-Forwarded-For", "X-Real-IP", "CF-Connecting-IP"}
	for _, h := range hdrs {
		if v := r.Header.Get(h); v != "" {
			parts := strings.Split(v, ",")
			candidate := strings.TrimSpace(parts[0])
			if candidate != "" {
				return candidate
			}
		}
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

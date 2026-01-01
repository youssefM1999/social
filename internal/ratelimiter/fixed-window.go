package ratelimiter

import (
	"sync"
	"time"
)

type FixedWindowLimiter struct {
	sync.RWMutex
	clients map[string]int
	limit   int
	window  time.Duration
}

func NewFixedWindowLimiter(limit int, window time.Duration) *FixedWindowLimiter {
	return &FixedWindowLimiter{
		clients: make(map[string]int),
		limit:   limit,
		window:  window,
	}
}

func (l *FixedWindowLimiter) Allow(ip string) (bool, time.Duration) {
	l.RLock()
	count, exists := l.clients[ip]
	l.RUnlock()

	if !exists || count < l.limit {
		l.Lock()
		if !exists {
			go l.resetCount(ip)
		}
		l.clients[ip]++
		l.Unlock()
		return true, 0

	}

	return false, l.window
}

func (l *FixedWindowLimiter) resetCount(ip string) {
	time.Sleep(l.window)
	l.Lock()
	delete(l.clients, ip)
	l.Unlock()
}

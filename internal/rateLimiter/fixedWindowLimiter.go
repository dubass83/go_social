package ratelimiter

import (
	"sync"
	"time"
)

type fixedWindowLimiter struct {
	clients map[string]int
	config  Config
	sync.RWMutex
}

func NewFixwedWindowLimeter(config Config) Limiter {
	return &fixedWindowLimiter{
		clients: make(map[string]int),
		config:  config,
	}
}

func (l *fixedWindowLimiter) Allow(ip string) (bool, time.Duration) {
	l.RLock()
	count, exist := l.clients[ip]
	l.RUnlock()
	// Allow request if IP is new or hasn't exceeded the limit
	// For new IPs: initialize counter, spawn goroutine to delete entry after time frame
	// For existing IPs: increment counter if below limit
	if !exist || count < l.config.RequestPerTimeFrame {
		l.Lock()
		if !exist {
			go l.resetCount(ip)
		}
		l.clients[ip]++
		l.Unlock()
		return true, 0
	}
	return false, l.config.TimeFrame
}

func (l *fixedWindowLimiter) resetCount(ip string) {
	time.Sleep(l.config.TimeFrame)
	l.Lock()
	delete(l.clients, ip)
	l.Unlock()
}

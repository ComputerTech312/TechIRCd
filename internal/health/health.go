package health

import (
	"log"
	"runtime"
	"sync/atomic"
	"time"
)

// HealthMonitor tracks server health metrics
type HealthMonitor struct {
	server        *Server
	totalClients  int64
	totalMessages int64
	startTime     time.Time
	ticker        *time.Ticker
	shutdown      chan bool
}

func NewHealthMonitor(server *Server) *HealthMonitor {
	return &HealthMonitor{
		server:    server,
		startTime: time.Now(),
		shutdown:  make(chan bool),
	}
}

func (h *HealthMonitor) Start() {
	h.ticker = time.NewTicker(5 * time.Minute) // Log stats every 5 minutes
	go h.monitor()
}

func (h *HealthMonitor) Stop() {
	if h.ticker != nil {
		h.ticker.Stop()
	}
	close(h.shutdown)
}

func (h *HealthMonitor) IncrementClients() {
	atomic.AddInt64(&h.totalClients, 1)
}

func (h *HealthMonitor) DecrementClients() {
	atomic.AddInt64(&h.totalClients, -1)
}

func (h *HealthMonitor) IncrementMessages() {
	atomic.AddInt64(&h.totalMessages, 1)
}

func (h *HealthMonitor) monitor() {
	for {
		select {
		case <-h.ticker.C:
			h.logHealthStats()
		case <-h.shutdown:
			return
		}
	}
}

func (h *HealthMonitor) logHealthStats() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	h.server.mu.RLock()
	clientCount := len(h.server.clients)
	channelCount := len(h.server.channels)
	h.server.mu.RUnlock()

	totalClients := atomic.LoadInt64(&h.totalClients)
	totalMessages := atomic.LoadInt64(&h.totalMessages)
	uptime := time.Since(h.startTime)

	log.Printf("Health Stats - Uptime: %v, Clients: %d, Channels: %d, Total Clients: %d, Total Messages: %d",
		uptime.Round(time.Second), clientCount, channelCount, totalClients, totalMessages)

	log.Printf("Memory Stats - Alloc: %d KB, Sys: %d KB, NumGC: %d, Goroutines: %d",
		bToKb(m.Alloc), bToKb(m.Sys), m.NumGC, runtime.NumGoroutine())

	// Alert if memory usage is high
	if m.Alloc > 100*1024*1024 { // 100MB
		log.Printf("WARNING: High memory usage detected: %d MB", bToMb(m.Alloc))
	}

	// Alert if too many goroutines
	if runtime.NumGoroutine() > 1000 {
		log.Printf("WARNING: High goroutine count: %d", runtime.NumGoroutine())
	}
}

func bToKb(b uint64) uint64 {
	return b / 1024
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

package services

import (
	"Hackathon/internal/core"
	"Hackathon/internal/core/structs"
	"fmt"
	"sync"
	"time"
)

type PollingService struct {
	service   core.PortStatisticsService
	portStats []structs.PortInfo
	mu        sync.RWMutex
}

func NewPollingService(service core.PortStatisticsService) *PollingService {
	return &PollingService{
		service:   service,
		portStats: []structs.PortInfo{},
	}
}

func (p *PollingService) StartPolling(interval time.Duration) {
	go func() {
		for {
			time.Sleep(interval)
			err := p.service.PollStatistics()
			if err != nil {
				fmt.Println("Error polling statistics:", err)
				continue
			}
			p.mu.Lock()
			p.portStats = p.service.GetPortStats()
			p.mu.Unlock()
		}
	}()
}

func (p *PollingService) GetPortStats() []structs.PortInfo {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.portStats
}

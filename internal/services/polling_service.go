package services

import (
	"Hackathon/internal/core"
	"Hackathon/internal/core/events"
	"Hackathon/internal/core/structs"
	"fmt"
	"sync"
	"time"
)

type PollingService struct {
	service   core.PortStatisticsService
	portStats []structs.PortInfo
	history   map[int][]structs.PortInfo
	mu        sync.RWMutex
}

func NewPollingService(service core.PortStatisticsService) *PollingService {
	service.FetchPorts()
	return &PollingService{
		service:   service,
		portStats: []structs.PortInfo{},
		history:   make(map[int][]structs.PortInfo),
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
			p.saveHistory()
			p.mu.Unlock()
		}
	}()
}

func (p *PollingService) saveHistory() {
	for _, port := range p.portStats {
		p.history[port.Index] = append(p.history[port.Index], port)

		if len(p.history[port.Index]) > 100 {
			p.history[port.Index] = p.history[port.Index][1:]
		}

		event := events.PortStatEvent{
			Port:    port,
			History: p.history[port.Index],
		}

		go events.GetDispatcher().Dispatch("PortStatUpdated", event)
	}
}

func (p *PollingService) GetPortStatus(index int) float64 {
	if !p.IsValidPortIndex(index) {
		return 0
	}
	if p.portStats[index].OperStatus == "UP" {
		return 1
	}
	return 0
}

func (p *PollingService) GetPortStats() []structs.PortInfo {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.portStats
}

func (p *PollingService) IsValidPortIndex(index int) bool {
	return index > 0 && index < len(p.portStats)
}

func (p *PollingService) GetHistoricStats(index int) []structs.PortInfo {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.history[index]
}

package services

import (
	"Hackathon/internal/core/snmp"
	"Hackathon/internal/core/structs"
	"fmt"
	"sync"
	"time"
)

type PollingService struct {
	snmpService *snmp.SnmpService
	portStats   []structs.PortInfo
	mu          sync.RWMutex
}

func NewPollingService(snmpService *snmp.SnmpService) *PollingService {
	return &PollingService{
		snmpService: snmpService,
		portStats:   []structs.PortInfo{},
	}
}

func (p *PollingService) StartPolling(interval time.Duration) {
	go func() {
		for {
			time.Sleep(interval)
			err := p.snmpService.PollStatistics()
			if err != nil {
				fmt.Println("Error polling statistics:", err)
				continue
			}
			p.mu.Lock()
			p.portStats = p.snmpService.PortStats
			p.mu.Unlock()
		}
	}()
}

func (p *PollingService) GetPortStats() []structs.PortInfo {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.portStats
}

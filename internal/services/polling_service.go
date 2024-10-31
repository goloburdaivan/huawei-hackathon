package services

import (
	"Hackathon/internal/core"
	"Hackathon/internal/core/structs"
	"errors"
	"fmt"
	"github.com/go-toast/toast"
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

			for _, port := range p.portStats {
				p.history[port.Index] = append(p.history[port.Index], port)

				if len(p.history[port.Index]) > 100 {
					p.history[port.Index] = p.history[port.Index][1:]
				}

				if port.InBandwidthActual > port.InBandwidthUtil || port.OutBandwidthActual > port.OutBandwidthUtil {
					message := fmt.Sprintf("ВНИМАНИЕ: Перегрузка на порте %d (%s).\nУтилизация входящей полосы: %.2f%%\nУтилизация исходящей полосы: %.2f%%",
						port.Index, port.Name, port.InBandwidthUtil, port.OutBandwidthUtil)
					p.sendToastNotification(message)
				}
			}
			p.mu.Unlock()
		}
	}()
}

func (p *PollingService) sendToastNotification(message string) {
	notification := toast.Notification{
		AppID:   "YourAppID",
		Title:   "Сетевое уведомление",
		Message: message,
	}
	err := notification.Push()
	if err != nil {
		fmt.Println("Ошибка отправки уведомления:", err)
	}
}

func (p *PollingService) PredictPortStat(index int) (structs.PortInfo, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	history, exists := p.history[index]
	if !exists || len(history) < 2 {
		return structs.PortInfo{}, errors.New("Недостаточно данных для прогнозирования")
	}

	n := float64(len(history))
	var sumX, sumY, sumXY, sumX2 float64

	for i, stat := range history {
		x := float64(i)
		y := float64(stat.InOctets)
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}

	denominator := n*sumX2 - sumX*sumX
	if denominator == 0 {
		return structs.PortInfo{}, errors.New("Ошибка вычисления регрессии")
	}
	a := (n*sumXY - sumX*sumY) / denominator
	b := (sumY - a*sumX) / n

	nextX := n
	predictedInOctets := a*nextX + b

	predictedStat := history[len(history)-1]
	predictedStat.InOctets = uint(predictedInOctets)

	return predictedStat, nil
}

func (p *PollingService) GetPortStats() []structs.PortInfo {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.portStats
}

func (p *PollingService) IsValidPortIndex(index int) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return index >= 0 && index < len(p.portStats)
}

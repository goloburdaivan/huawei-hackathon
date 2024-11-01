package listeners

import (
	"Hackathon/internal/core/events"
	"Hackathon/internal/services"
	"fmt"
	"os"
)

type BandwidthCriticalListener struct{}

func (n *BandwidthCriticalListener) Handle(e events.Event) {
	if os.Getenv("NOTIFICATION_ON") == "false" {
		return
	}

	if event, ok := e.(events.PortStatEvent); ok {
		port := event.Port
		if port.InBandwidthActual > port.InBandwidthUtil ||
			port.OutBandwidthActual > port.OutBandwidthUtil &&
				(port.OperStatus != "DOWN" || port.AdminStatus != "DOWN") {
			message := fmt.Sprintf(
				"ВНИМАНИЕ: Перегрузка на порте %d (%s).\nУтилизация входящей полосы: %.2f%%\nУтилизация исходящей полосы: %.2f%%",
				port.Index, port.Name, port.InBandwidthActual, port.OutBandwidthActual,
			)
			services.SendToastNotification(message)
		}
	}
}

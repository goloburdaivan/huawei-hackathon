package listeners

import (
	"Hackathon/internal/core/events"
	"Hackathon/internal/services"
	"fmt"
	"os"
)

type PortStatusListener struct{}

func (n *PortStatusListener) Handle(e events.Event) {
	if os.Getenv("NOTIFICATION_ON") == "false" {
		return
	}

	if event, ok := e.(events.PortStatEvent); ok {
		if len(event.History) < 2 {
			return
		}

		lastStatus := event.History[len(event.History)-2].OperStatus
		currentStatus := event.Port.OperStatus

		if lastStatus != currentStatus {
			services.SendToastNotification(fmt.Sprintf("WARNING: Port %s changed status to %s", event.Port.Name, currentStatus))
		}
	}
}

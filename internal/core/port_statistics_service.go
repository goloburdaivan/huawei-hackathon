package core

import "Hackathon/internal/core/structs"

type PortStatisticsService interface {
	PollStatistics() error
	GetPortStats() []structs.PortInfo
	Connect() error
	CloseConnection()
	FetchPorts() error
}

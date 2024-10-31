package events

import "Hackathon/internal/core/structs"

type PortStatEvent struct {
	Port    structs.PortInfo
	History []structs.PortInfo
}

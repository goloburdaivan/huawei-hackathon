package snmp

import "strings"

func getStatusLabel(status int) string {
	switch status {
	case 1:
		return "UP"
	case 2:
		return "DOWN"
	default:
		return "UNKNOWN"
	}
}

func GetPortStatus(operStatus string) float64 {
	if strings.ToLower(operStatus) == "up" {
		return 1
	}
	return 0
}

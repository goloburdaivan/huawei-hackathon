package snmp

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
	if operStatus == "UP" {
		return 1
	}
	return 0
}

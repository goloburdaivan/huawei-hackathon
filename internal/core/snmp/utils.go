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

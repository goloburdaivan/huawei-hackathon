package snmp

import (
	"fmt"
	"regexp"
)

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

func isValidIP(ip string) bool {
	ipRegex := regexp.MustCompile(`^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$`)
	if !ipRegex.MatchString(ip) {
		return false
	}

	var a, b, c, d int
	_, err := fmt.Sscanf(ip, "%d.%d.%d.%d", &a, &b, &c, &d)
	if err != nil || a > 255 || b > 255 || c > 255 || d > 255 {
		return false
	}
	return true
}

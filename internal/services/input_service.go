package services

import (
	"fmt"
	"regexp"
)

func IsValidIP(ip string) bool {
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

func GetSSHInput() (ip string, port uint16, user, password string) {
	for {
		fmt.Print("Enter the IP address for SSH (example: 192.168.10.2): ")
		fmt.Scanln(&ip)
		if IsValidIP(ip) {
			break
		}
		fmt.Println("Invalid IP address format. Please try again.")
	}

	fmt.Print("Enter the port for SSH (22): ")
	fmt.Scanln(&port)
	fmt.Print("Enter the username for SSH (admin): ")
	fmt.Scanln(&user)
	fmt.Print("Enter the password for SSH (admin123): ")
	fmt.Scanln(&password)
	return
}

func GetSNMPInput() (ip string, port uint16, community string) {
	for {
		fmt.Print("Enter the IP address for SNMP (example: 192.168.10.2): ")
		fmt.Scanln(&ip)
		if IsValidIP(ip) {
			break
		}
		fmt.Println("Invalid IP address format. Please try again.")
	}

	fmt.Print("Enter the port for SNMP (example: 161): ")
	fmt.Scanln(&port)
	fmt.Print("Enter the community string for SNMP (public): ")
	fmt.Scanln(&community)
	return
}

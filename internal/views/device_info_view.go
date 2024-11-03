package views

import (
	"Hackathon/internal/core/structs"
	"fmt"
)

func DisplayDeviceInfo(deviceStatus *structs.DeviceStatus) {
	fmt.Println("Device Status Information:")
	fmt.Printf("Slot: %d\n", deviceStatus.Slot)
	fmt.Printf("Subslot: %s\n", deviceStatus.Sub)
	fmt.Printf("Type: %s\n", deviceStatus.Type)
	fmt.Printf("Online: %s\n", deviceStatus.Online)
	fmt.Printf("Power: %s\n", deviceStatus.Power)
	fmt.Printf("Registration: %s\n", deviceStatus.Register)
	fmt.Printf("Status: %s\n", deviceStatus.Status)
	fmt.Printf("Role: %s\n", deviceStatus.Role)
}

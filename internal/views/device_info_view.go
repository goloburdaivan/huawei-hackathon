package views

import (
	"Hackathon/internal/core/structs"
	"fmt"
)

func DisplayDeviceInfo(deviceStatus *structs.DeviceStatus) {
	fmt.Println("Информация о статусе устройства:")
	fmt.Printf("Слот: %d\n", deviceStatus.Slot)
	fmt.Printf("Подслот: %s\n", deviceStatus.Sub)
	fmt.Printf("Тип: %s\n", deviceStatus.Type)
	fmt.Printf("Онлайн: %s\n", deviceStatus.Online)
	fmt.Printf("Мощность: %s\n", deviceStatus.Power)
	fmt.Printf("Регистрация: %s\n", deviceStatus.Register)
	fmt.Printf("Статус: %s\n", deviceStatus.Status)
	fmt.Printf("Роль: %s\n", deviceStatus.Role)
}

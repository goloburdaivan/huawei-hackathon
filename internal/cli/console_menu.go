package cli

import (
	"Hackathon/internal/controllers"
	"fmt"
)

type ConsoleMenu struct {
	portController *controllers.PortController
}

func NewConsoleMenu(portController *controllers.PortController) *ConsoleMenu {
	return &ConsoleMenu{
		portController: portController,
	}
}

func (m *ConsoleMenu) DisplayMenu() {
	for {
		fmt.Println("1. Показать информацию о портах")
		fmt.Println("2. Выйти")
		var choice int
		fmt.Scanln(&choice)
		switch choice {
		case 1:
			m.portController.ShowPortStats()
		case 2:
			fmt.Println("Завершение работы...")
			return
		default:
			fmt.Println("Неверный выбор. Попробуйте снова.")
		}
	}
}

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
		fmt.Println("1. Show information about ports")
		fmt.Println("2. Exit")
		var choice int
		fmt.Scanln(&choice)
		switch choice {
		case 1:
			m.portController.ShowPortStats()
		case 2:
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}

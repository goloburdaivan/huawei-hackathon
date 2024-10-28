package controllers

import (
	"Hackathon/internal/services"
	"Hackathon/internal/views"
	"fmt"
)

type PortController struct {
	pollingService *services.PollingService
	stopChannel    chan bool
}

func NewPortController(pollingService *services.PollingService) *PortController {
	return &PortController{
		pollingService: pollingService,
		stopChannel:    make(chan bool),
	}
}

func (pc *PortController) ShowPortStats() {
	go func() {
		fmt.Println("Нажмите Enter, чтобы вернуться в меню.")
		fmt.Scanln()
		pc.stopChannel <- true
	}()

	for {
		select {
		case <-pc.stopChannel:
			fmt.Println("Возвращаемся в меню...")
			return
		default:
			portStats := pc.pollingService.GetPortStats()
			views.DisplayPortStats(portStats)
		}
	}
}

func (pc *PortController) ShowPortGraph() {
	portIndex := pc.promptPortIndex()

	if portIndex == -1 {
		fmt.Println("Возвращаемся в меню...")
		return
	}

	portStats := pc.pollingService.GetPortStats()
	if portIndex < 0 || portIndex >= len(portStats) {
		fmt.Printf("Порт с индексом %d не существует.\n", portIndex)
		return
	}

	portName := portStats[portIndex].Name

	// Функция для получения статуса порта (0 или 1)
	getStatus := func() float64 {
		portStats := pc.pollingService.GetPortStats()
		if portStats[portIndex].OperStatus == "UP" {
			return 1
		}
		return 0
	}

	views.DisplayPortGraph(portName, portIndex, getStatus, pc.stopChannel)
	fmt.Println("Возвращаемся в меню...")
}

func (pc *PortController) promptPortIndex() int {
	for {
		var portIndex int
		fmt.Println("Введите индекс порта для отображения графика (введите -1 для возврата в меню):")

		_, err := fmt.Scanln(&portIndex)
		if err != nil {
			fmt.Println("Ошибка ввода, пожалуйста, введите корректный индекс или -1 для возврата в меню.")
			continue
		}

		if portIndex == -1 {
			return -1
		}

		portStats := pc.pollingService.GetPortStats()
		if portIndex >= 0 && portIndex < len(portStats) {
			return portIndex
		}

		fmt.Printf("Порт с индексом %d не найден. Попробуйте снова или введите -1 для возврата в меню.\n", portIndex)
	}
}

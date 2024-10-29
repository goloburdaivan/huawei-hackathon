package controllers

import (
	"Hackathon/internal/services"
	"Hackathon/internal/views"
	"errors"
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
	portIndex, err := pc.getPortIndex()
	if err != nil {
		fmt.Println(err)
		fmt.Println("Возвращаемся в меню...")
		return
	}

	portName, err := pc.pollingService.GetPortName(portIndex)
	if err != nil {
		fmt.Println(err)
		return
	}
	getStatus := func() float64 {
		return pc.pollingService.GetPortStatus(portIndex)
	}

	views.DisplayPortGraph(portName, portIndex, getStatus, pc.stopChannel)
	fmt.Println("Возвращаемся в меню...")
}

func (pc *PortController) getPortIndex() (int, error) {
	for {
		var portIndex int
		fmt.Println("Введите индекс порта для отображения графика (введите -1 для возврата в меню):")

		_, err := fmt.Scanln(&portIndex)
		if err != nil {
			fmt.Println("Ошибка ввода, пожалуйста, введите корректный индекс или -1 для возврата в меню.")
			continue
		}

		if portIndex == -1 {
			return -1, errors.New("отмена: возврат в меню")
		}

		if pc.pollingService.IsValidPortIndex(portIndex) {
			return portIndex, nil
		}

		fmt.Printf("Порт с индексом %d не найден. Попробуйте снова или введите -1 для возврата в меню.\n", portIndex)
	}
}

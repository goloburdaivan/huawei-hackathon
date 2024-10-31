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

func (pc *PortController) ShowPortStatusGraph() {
	portStats := pc.pollingService.GetPortStats()
	views.DisplayPortList(portStats)

	portIndex, err := pc.getPortIndex()
	if err != nil {
		fmt.Println(err)
		fmt.Println("Возвращаемся в меню...")
		return
	}
	portName := portStats[portIndex].Name

	views.DisplayPortStatusGraph(portName, portIndex, pc.pollingService, pc.stopChannel)
	fmt.Println("Возвращаемся в меню...")
}

func (pc *PortController) ShowPortOctetsGraph(octetType string) {
	portStats := pc.pollingService.GetPortStats()
	views.DisplayPortList(portStats)

	portIndex, err := pc.getPortIndex()
	if err != nil {
		fmt.Println(err)
		fmt.Println("Возвращаемся в меню...")
		return
	}

	portName := portStats[portIndex].Name

	var getOctetsFunc func() float64
	if octetType == "InOctets" {
		getOctetsFunc = func() float64 {
			return float64(pc.pollingService.GetPortStats()[portIndex].InOctets)
		}
	} else if octetType == "OutOctets" {
		getOctetsFunc = func() float64 {
			return float64(pc.pollingService.GetPortStats()[portIndex].OutOctets)
		}
	} else {
		fmt.Println("Неизвестный тип Octets. Используйте 'InOctets' или 'OutOctets'.")
		return
	}

	views.DisplayPortOctetsGraph(portName, portIndex, octetType, pc.stopChannel, getOctetsFunc)
	fmt.Println("Возвращаемся в меню...")
}

func (pc *PortController) ShowPort() {
	portIndex, err := pc.getPortIndex()
	if err != nil {
		fmt.Println(err)
		fmt.Println("Возвращаемся в меню...")
		return
	}
	portStat := pc.pollingService.GetPortStats()
	views.DisplaySinglePortStats(&portStat[portIndex])

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
			return portIndex - 1, nil
		}

		fmt.Printf("Порт с индексом %d не найден. Попробуйте снова или введите -1 для возврата в меню.\n", portIndex)
	}
}

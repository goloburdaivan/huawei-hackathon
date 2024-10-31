package controllers

import (
	"Hackathon/internal/services"
	"Hackathon/internal/views"
	"errors"
	"fmt"
)

type PortController struct {
	pollingService    *services.PollingService
	stopChannel       chan bool
	predictionService *services.PredictionService
}

func NewPortController(pollingService *services.PollingService) *PortController {
	return &PortController{
		pollingService:    pollingService,
		stopChannel:       make(chan bool),
		predictionService: services.NewPredictionService(),
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

	portStats := pc.pollingService.GetPortStats()
	portName := portStats[portIndex].Name
	portStatus := portStats[portIndex].OperStatus

	views.DisplayPortGraph(portName, portIndex, portStatus, pc.stopChannel)

	fmt.Println("Возвращаемся в меню...")
}

func (pc *PortController) ShowPort() {
	portIndex, err := pc.getPortIndex()
	if err != nil {
		fmt.Println(err)
		fmt.Println("Возвращаемся в меню...")
		return
	}
	portStat := pc.pollingService.GetPortStats()[portIndex]
	views.DisplaySinglePortStats(&portStat)

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

func (c *PortController) ShowPortPrediction() {
	fmt.Println("Введите индекс порта для прогнозирования:")
	var index int
	fmt.Scanln(&index)

	if !c.pollingService.IsValidPortIndex(index) {
		fmt.Println("Неверный индекс порта.")
		return
	}

	predictedStat, err := c.predictionService.PredictPortStat(c.pollingService.GetHistoricStats(index))
	if err != nil {
		fmt.Println("Ошибка прогнозирования:", err)
		return
	}

	fmt.Printf("Прогнозируемые данные для порта %d:\n", index)
	fmt.Printf("InOctets: %d\n", predictedStat.InOctets)
	fmt.Printf("OutOctets: %d\n", predictedStat.OutOctets)
}

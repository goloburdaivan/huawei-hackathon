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
		fmt.Println("Press Enter to return to the menu.")
		fmt.Scanln()
		pc.stopChannel <- true
	}()

	for {
		select {
		case <-pc.stopChannel:
			fmt.Println("Returning to the menu...")
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
		fmt.Println("Returning to the menu...")
		return
	}
	portName := portStats[portIndex].Name

	views.DisplayPortStatusGraph(portName, portIndex, pc.pollingService, pc.stopChannel)
	fmt.Println("Returning to the menu...")
}

func (pc *PortController) ShowPortGrowthGraph(growthType string) {
	portStats := pc.pollingService.GetPortStats()
	views.DisplayPortList(portStats)

	portIndex, err := pc.getPortIndex()
	if err != nil {
		fmt.Println(err)
		fmt.Println("Returning to the menu...")
		return
	}

	portName := portStats[portIndex].Name

	var getGrowthFunc func() float64

	switch growthType {
	case "InOctets":
		getGrowthFunc = func() float64 {
			return float64(pc.pollingService.GetPortStats()[portIndex].InOctets)
		}
	case "OutOctets":
		getGrowthFunc = func() float64 {
			return float64(pc.pollingService.GetPortStats()[portIndex].OutOctets)
		}
	case "InBandwidth":
		getGrowthFunc = func() float64 {
			return float64(pc.pollingService.GetPortStats()[portIndex].InBandwidthActual)
		}
	case "OutBandwidth":
		getGrowthFunc = func() float64 {
			return float64(pc.pollingService.GetPortStats()[portIndex].OutBandwidthActual)
		}
	default:
		fmt.Println("Unknown type. Use 'InOctets', 'OutOctets', 'InBandwidth', or 'OutBandwidth'.")
		return
	}

	views.DisplayPortGrowthGraph(portName, portIndex, growthType, pc.stopChannel, getGrowthFunc)
	fmt.Println("Returning to the menu...")
}

func (pc *PortController) ShowPort() {
	portIndex, err := pc.getPortIndex()
	if err != nil {
		fmt.Println(err)
		fmt.Println("Returning to the menu...")
		return
	}
	portStat := pc.pollingService.GetPortStats()
	views.DisplaySinglePortStats(&portStat[portIndex])

	fmt.Println("Returning to the menu...")
}

func (pc *PortController) getPortIndex() (int, error) {
	for {
		var portIndex int
		fmt.Println("Enter the port index for the graph display (enter -1 to return to the menu):")

		_, err := fmt.Scanln(&portIndex)
		if err != nil {
			fmt.Println("Input error, please enter a valid index or -1 to return to the menu.")
			continue
		}

		if portIndex == -1 {
			return -1, errors.New("canceled: returning to the menu")
		}

		if pc.pollingService.IsValidPortIndex(portIndex) {
			return portIndex - 1, nil
		}

		fmt.Printf("Port with index %d not found. Please try again or enter -1 to return to the menu.\n", portIndex)
	}
}

func (c *PortController) ShowPortPrediction() {
	fmt.Println("Enter the port index for prediction:")
	var index int
	fmt.Scanln(&index)

	if !c.pollingService.IsValidPortIndex(index) {
		fmt.Println("Invalid port index.")
		return
	}

	predictedStat, err := c.predictionService.PredictPortStat(c.pollingService.GetHistoricStats(index))
	if err != nil {
		fmt.Println("Prediction error:", err)
		return
	}

	fmt.Printf("Predicted data for port %d:\n", index)
	fmt.Printf("InOctets: %d\n", predictedStat.InOctets)
	fmt.Printf("OutOctets: %d\n", predictedStat.OutOctets)
}

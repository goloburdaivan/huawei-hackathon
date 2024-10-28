package controllers

import (
	"Hackathon/internal/services"
	"Hackathon/internal/views"
	"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"time"
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
	if err := ui.Init(); err != nil {
		fmt.Printf("failed to initialize termui: %v\n", err)
		return
	}
	defer ui.Close()

	plot := widgets.NewPlot()
	plot.Title = "Port Status (UP = 1, DOWN = 0)"
	plot.Data = [][]float64{{}}
	plot.SetRect(0, 0, 110, 10)
	plot.AxesColor = ui.ColorWhite
	plot.Marker = widgets.MarkerBraille
	plot.HorizontalScale = 1

	plot.Data[0] = append(plot.Data[0], 0)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	uiEvents := ui.PollEvents()

	go func() {
		for {
			select {
			case <-ticker.C:
				var status float64
				portStats := pc.pollingService.GetPortStats()[6]
				if portStats.OperStatus == "UP" {
					status = 1
					plot.LineColors[0] = ui.ColorGreen
				} else {
					status = 0
					plot.LineColors[0] = ui.ColorRed
				}

				if len(plot.Data[0]) > 100 {
					plot.Data[0] = plot.Data[0][1:]
				}

				plot.Data[0] = append(plot.Data[0], status)

				currentTime := time.Now().Format("15:04:05")
				plot.Title = fmt.Sprintf("Port Status (UP = 1, DOWN = 0) - Time: %s", currentTime)

				ui.Render(plot)

			case e := <-uiEvents:
				if e.ID == "q" || e.ID == "<Enter>" {
					pc.stopChannel <- true
					return
				}
			}
		}
	}()

	for {
		select {
		case <-pc.stopChannel:
			fmt.Println("Возвращаемся в меню...")
			return
		}
	}
}

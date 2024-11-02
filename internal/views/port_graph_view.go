package views

import (
	"Hackathon/internal/services"
	"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"time"
)

func DisplayPortStatusGraph(portName string, portIndex int, pollingService *services.PollingService, stopChannel chan bool) {
	if err := ui.Init(); err != nil {
		fmt.Printf("Failed to initialize termui: %v\n", err)
		return
	}
	defer ui.Close()

	baseTitle := fmt.Sprintf("Port %s (Index: %d) Status (UP = 1, DOWN = 0)", portName, portIndex+1)
	plot := initializePlot(baseTitle, 10)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	uiEvents := ui.PollEvents()

	go func() {
		for {
			select {
			case <-ticker.C:
				currentTime := time.Now().Format("15:04:05")
				updateStatusPlotData(plot, pollingService.GetPortStatus(portIndex), baseTitle, currentTime)
				ui.Render(plot)

			case e := <-uiEvents:
				if e.ID == "q" || e.ID == "<Enter>" {
					stopChannel <- true
					return
				}
			}
		}
	}()

	<-stopChannel
	clearConsole()
}

func DisplayPortGrowthGraph(portName string, portIndex int, growthType string, stopChannel chan bool, getGrowth func() float64) {
	if err := ui.Init(); err != nil {
		fmt.Printf("Failed to initialize termui: %v\n", err)
		return
	}
	defer ui.Close()

	baseTitle := fmt.Sprintf("Port %s (Index: %d) %s Over Time", portName, portIndex+1, growthType)
	plot := initializePlot(baseTitle, 20)
	plot.Data[0] = append(plot.Data[0], 0)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	uiEvents := ui.PollEvents()

	go func() {
		for {
			select {
			case <-ticker.C:
				currentTime := time.Now().Format("15:04:05")
				growth := getGrowth()

				if growthType == "InOctets" || growthType == "OutOctets" {
					growth = growth / 1024 / 1024
				}
				updateOctetsPlotData(plot, growth, baseTitle, currentTime)
				ui.Render(plot)

			case e := <-uiEvents:
				if e.ID == "q" || e.ID == "<Enter>" {
					stopChannel <- true
					return
				}
			}
		}
	}()

	<-stopChannel
	clearConsole()
}

func initializePlot(baseTitle string, y2 int) *widgets.Plot {
	plot := widgets.NewPlot()
	plot.Title = baseTitle
	plot.Data = [][]float64{{0}}
	plot.SetRect(0, 0, 110, y2)
	plot.AxesColor = ui.ColorWhite
	plot.Marker = widgets.MarkerBraille
	plot.HorizontalScale = 1
	return plot
}

func updateOctetsPlotData(plot *widgets.Plot, octets float64, baseTitle, currentTime string) {
	if len(plot.Data) == 0 {
		plot.Data = append(plot.Data, []float64{})
	}

	if len(plot.Data[0]) > 100 {
		plot.Data[0] = plot.Data[0][1:]
	}

	plot.Data[0] = append(plot.Data[0], octets)
	plot.Title = fmt.Sprintf("%s - Time: %s", baseTitle, currentTime)
	plot.LineColors = []ui.Color{ui.ColorCyan}
}

func updateStatusPlotData(plot *widgets.Plot, status float64, baseTitle, currentTime string) {
	if len(plot.Data) == 0 {
		plot.Data = append(plot.Data, []float64{})
	}

	if len(plot.Data[0]) > 100 {
		plot.Data[0] = plot.Data[0][1:]
	}

	plot.Data[0] = append(plot.Data[0], status)
	plot.Title = fmt.Sprintf("%s - Time: %s", baseTitle, currentTime)

	if status == 1 {
		plot.LineColors = []ui.Color{ui.ColorGreen}
	} else {
		plot.LineColors = []ui.Color{ui.ColorRed}
	}
}

package views

import (
	"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"time"
)

func DisplayPortGraph(portName string, portIndex int, getStatus func() float64, stopChannel chan bool) {
	if err := ui.Init(); err != nil {
		fmt.Printf("Failed to initialize termui: %v\n", err)
		return
	}
	defer ui.Close()

	baseTitle := fmt.Sprintf("Port %s (Index: %d) Status (UP = 1, DOWN = 0)", portName, portIndex)
	plot := initializePlot(baseTitle)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	uiEvents := ui.PollEvents()

	go func() {
		for {
			select {
			case <-ticker.C:
				currentTime := time.Now().Format("15:04:05")
				updatePlotData(plot, getStatus(), baseTitle, currentTime)
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
	ClearConsole()
}

func initializePlot(baseTitle string) *widgets.Plot {
	plot := widgets.NewPlot()
	plot.Title = baseTitle
	plot.Data = [][]float64{{0}}
	plot.SetRect(0, 0, 110, 10)
	plot.AxesColor = ui.ColorWhite
	plot.Marker = widgets.MarkerBraille
	plot.HorizontalScale = 1
	return plot
}

func updatePlotData(plot *widgets.Plot, status float64, baseTitle, currentTime string) {
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

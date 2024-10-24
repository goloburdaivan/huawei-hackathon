package controllers

import (
	"Hackathon/internal/services"
	"fmt"
)

type ExportController struct {
	exportService  *services.ExportService
	pollingService *services.PollingService
}

func NewExportController(exportService *services.ExportService, pollingService *services.PollingService) *ExportController {
	return &ExportController{
		exportService:  exportService,
		pollingService: pollingService,
	}
}

func (ec *ExportController) ExportPortStats() {
	portStats := ec.pollingService.GetPortStats()
	go ec.exportService.ExportPortStatsToCSV(portStats)
}

func (ec *ExportController) ExportPortStatsByPort() {
	var portName string
	fmt.Print("Введите название порта для экспорта: ")
	fmt.Scanln(&portName)

	go func() {
		portStats := ec.pollingService.GetPortStats()
		ec.exportService.ExportPortStatsByPort(portStats, portName)
	}()
}

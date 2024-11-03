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
	var portIndex int
	fmt.Print("Enter the index of the port to export: ")
	fmt.Scanln(&portIndex)

	go func() {
		portStats := ec.pollingService.GetPortStats()
		ec.exportService.ExportPortStatsByIndex(portStats, portIndex)
	}()
}

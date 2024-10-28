package controllers

import (
	"Hackathon/internal/services"
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

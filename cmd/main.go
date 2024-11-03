package main

import (
	"Hackathon/internal/cli"
	"Hackathon/internal/controllers"
	"Hackathon/internal/core/events"
	"Hackathon/internal/core/events/listeners"
	"Hackathon/internal/services"
	"github.com/joho/godotenv"
	"log"
	"time"
)

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	dispatcher := events.GetDispatcher()
	dispatcher.Register("PortStatUpdated", &listeners.BandwidthCriticalListener{})
	dispatcher.Register("PortStatUpdated", &listeners.PortStatusListener{})

	sshService := services.ConnectSSH()
	pollingService := services.NewPollingService(sshService)
	pollingService.StartPolling(1 * time.Second)
	exportService := services.NewExportService()
	portController := controllers.NewPortController(pollingService)
	exportController := controllers.NewExportController(exportService, pollingService)
	sshController := controllers.NewSSHController(sshService)

	menu := cli.NewMenuBuilder("Main Menu").
		AddAction("Show information about all ports", portController.ShowPortStats).
		AddSubMenu("Export port information").
		AddAction("Export information about all ports", exportController.ExportPortStats).
		AddAction("Export information for a specific port", exportController.ExportPortStatsByPort).
		EndSubMenu().
		AddAction("Show status graph for a specific port", portController.ShowPortStatusGraph).
		AddSubMenu("Graphs for InOctets/OutOctets for a specific port").
		AddAction("Show InOctets graph for a specific port", func() { portController.ShowPortGrowthGraph("InOctets") }).
		AddAction("Show OutOctets graph for a specific port", func() { portController.ShowPortGrowthGraph("OutOctets") }).
		EndSubMenu().
		AddSubMenu("Graphs for InBandwidth/OutBandwidth for a specific port").
		AddAction("Show InBandwidth graph for a specific port", func() { portController.ShowPortGrowthGraph("InBandwidth") }).
		AddAction("Show OutBandwidth graph for a specific port", func() { portController.ShowPortGrowthGraph("OutBandwidth") }).
		EndSubMenu().
		AddAction("Display information for a specific port", portController.ShowPort).
		AddAction("Predict port statistics", portController.ShowPortPrediction).
		EndSubMenu().
		AddAction("Start SSH CLI session", sshController.StartCliSession).
		Build()

	menu.Execute()
}

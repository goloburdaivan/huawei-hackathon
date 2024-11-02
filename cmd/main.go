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

	menu := cli.NewMenuBuilder("Главное меню").
		AddAction("Показать информацию о всех портах", portController.ShowPortStats).
		AddSubMenu("Экспорт информации про порты").
		AddAction("Экспортировать информацию о всех портах", exportController.ExportPortStats).
		AddAction("Экспортировать информацию конкретного порта", exportController.ExportPortStatsByPort).
		EndSubMenu().
		AddAction("Показать график статусов для определённого порта", portController.ShowPortStatusGraph).
		AddSubMenu("Графики InOctets/OutOctets для определённого порта").
		AddAction("Показать график InOctets для определённого порта", func() { portController.ShowPortOctetsGraph("InOctets") }).
		AddAction("Показать график OutOctets для определённого порта", func() { portController.ShowPortOctetsGraph("OutOctets") }).
		EndSubMenu().
		AddAction("Вывести информацию по определённому порту", portController.ShowPort).
		AddAction("Прогнозировать статистику порта", portController.ShowPortPrediction).
		EndSubMenu().
		AddAction("Начать SSH CLI сессию", sshController.StartCliSession).
		Build()

	menu.Execute()
}

package main

import (
	"Hackathon/internal/cli"
	"Hackathon/internal/controllers"
	"Hackathon/internal/core/events"
	"Hackathon/internal/core/events/listeners"
	"Hackathon/internal/services"
	"fmt"
	"time"
)

func main() {
	sshService := services.ConnectSSH()
	snmpService := services.ConnectSNMP()
	defer snmpService.CloseConnection()
	err := snmpService.Connect()
	dispatcher := events.GetDispatcher()
	dispatcher.Register("PortStatUpdated", &listeners.BandwidthCriticalListener{})
	err = sshService.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}

	pollingService := services.NewPollingService(snmpService)
	pollingService.StartPolling(1 * time.Second)
	exportService := services.NewExportService()

	portController := controllers.NewPortController(pollingService)
	exportController := controllers.NewExportController(exportService, pollingService)

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
		Build()

	menu.Execute()
}

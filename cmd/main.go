package main

import (
	"Hackathon/internal/cli"
	"Hackathon/internal/controllers"
	"Hackathon/internal/services"
	"fmt"
	"time"
)

func main() {
	sshService := services.ConnectSSH()
	snmpService := services.ConnectSNMP()
	defer snmpService.CloseConnection()

	err := snmpService.FetchPorts()
	if err != nil {
		fmt.Println("Ошибка получения данных о портах:", err)
		return
	}

	pollingService := services.NewPollingService(snmpService)
	pollingService.StartPolling(1 * time.Second)
	exportService := services.NewExportService()

	portController := controllers.NewPortController(pollingService)
	exportController := controllers.NewExportController(exportService, pollingService)
	deviceController := controllers.NewDeviceController(sshService)

	menu := cli.NewMenuBuilder("Главное меню").
		AddAction("Показать информацию о всех портах", portController.ShowPortStats).
		AddSubMenu("Экспорт информации про порты").
		AddAction("Экспортировать информацию о всех портах", exportController.ExportPortStats).
		AddAction("Экспортировать информацию конкретного порта", exportController.ExportPortStatsByPort).
		EndSubMenu().
		AddAction("Показать график для определённого порта", portController.ShowPortGraph).
		AddAction("Вывести информацию по определённому порту", portController.ShowPort).
		AddAction("Вывести информацию об устройстве", deviceController.ShowDeviceInfo).
		Build()

	menu.Execute()
}

package main

import (
	"Hackathon/internal/cli"
	"Hackathon/internal/controllers"
	"Hackathon/internal/services"
	"fmt"
	"time"
)

func main() {
	// sshService := services.ConnectSSH()
	snmpService := services.ConnectSNMP()
	defer snmpService.CloseConnection()
	err := snmpService.Connect()
	if err != nil {
		fmt.Println("Ошибка подключения к SNMP:", err)
		return
	}

	err = snmpService.FetchPorts()
	if err != nil {
		fmt.Println("Ошибка получения данных о портах:", err)
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
		Build()

	menu.Execute()
}

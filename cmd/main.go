package main

import (
	"Hackathon/internal/cli"
	"Hackathon/internal/controllers"
	"Hackathon/internal/core/snmp"
	"Hackathon/internal/core/ssh"
	"Hackathon/internal/services"
	"fmt"
	"time"
)

func main() {
	sshService := ssh.NewSshService("192.168.10.2", 22, "admin", "admin123")
	err := sshService.Connect()

	snmpService := snmp.NewSnmpService("192.168.10.2", 161, "public")
	defer snmpService.CloseConnection()
	err = snmpService.Connect()
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
		AddAction("Показать график для определённого порта", portController.ShowPortGraph).
		AddAction("Вывести информацию по определённому порту", portController.ShowPort).
		Build()

	menu.Execute()
}

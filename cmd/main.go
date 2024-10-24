package main

import (
	"Hackathon/internal/cli"
	"Hackathon/internal/controllers"
	"Hackathon/internal/core/snmp"
	"Hackathon/internal/services"
	"fmt"
	"time"
)

func main() {
	snmpService := snmp.NewSnmpService("192.168.10.2", 161, "public")
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
	sshController := controllers.NewSSHController("192.168.10.2", "22", "username", "password")

	menu := cli.NewMenuBuilder("Главное меню").
		AddAction("Показать информацию о портах", portController.ShowPortStats).
		AddAction("Экспортировать статистику портов в CSV", exportController.ExportPortStats).
		AddAction("Открыть SSH-сессию", sshController.StartSSHSession).
		Build()

	menu.Execute()
}

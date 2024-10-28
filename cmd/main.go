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
		AddAction("Показать информацию о портах", portController.ShowPortStats).
		AddAction("Показать график для портов", portController.ShowPortGraph).
		AddAction("Экспортировать статистику портов в CSV", exportController.ExportPortStats).
		Build()

	menu.Execute()
}

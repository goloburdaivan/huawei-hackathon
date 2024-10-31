package main

import (
	"Hackathon/internal/cli"
	"Hackathon/internal/controllers"
	"Hackathon/internal/core/events"
	"Hackathon/internal/core/events/listeners"
	"Hackathon/internal/core/snmp"
	"Hackathon/internal/core/ssh"
	"Hackathon/internal/services"
	"fmt"
	"time"
)

func main() {
	dispatcher := events.GetDispatcher()
	dispatcher.Register("PortStatUpdated", &listeners.BandwidthCriticalListener{})
	sshService := ssh.NewSshService("192.168.65.6", 22, "Student_1", "UY2AEaZ7BmKs#")
	err := sshService.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}

	snmpService := snmp.NewSnmpService("192.168.65.6", 161, "public")
	defer snmpService.CloseConnection()
	err = snmpService.Connect()
	if err != nil {
		fmt.Println("Ошибка подключения к SNMP:", err)
		return
	}

	pollingService := services.NewPollingService(sshService)
	pollingService.StartPolling(1 * time.Second)
	exportService := services.NewExportService()

	portController := controllers.NewPortController(pollingService)
	exportController := controllers.NewExportController(exportService, pollingService)

	menu := cli.NewMenuBuilder("Главное меню").
		AddAction("Показать информацию о портах", portController.ShowPortStats).
		AddSubMenu("Выберите какую информацию вы хотите вывести:").
		AddAction("Экспортировать информацию о всех портах", exportController.ExportPortStats).
		AddAction("Экспортировать информацию конкретного порта", exportController.ExportPortStatsByPort).
		EndSubMenu().
		AddAction("Показать график для портов", portController.ShowPortGraph).
		AddAction("Вывести информацию по определённому порту", portController.ShowPort).
		AddAction("Прогнозировать статистику порта", portController.ShowPortPrediction).
		Build()

	menu.Execute()
}

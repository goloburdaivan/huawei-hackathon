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

	portController := controllers.NewPortController(pollingService)

	menu := cli.NewMenuBuilder("Главное меню").
		AddAction("Показать информацию о портах", portController.ShowPortStats).
		AddAction("Показать график для портов", portController.ShowPortGraph).
		Build()

	menu.Execute()
}

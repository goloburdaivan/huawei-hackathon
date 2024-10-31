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
	var sshService *ssh.SshService
	for {
		sshIP, sshPort, sshUser, sshPass := services.GetSSHInput() // Используем функцию из input_service
		sshService = ssh.NewSshService(sshIP, sshPort, sshUser, sshPass)
		err := sshService.Connect()
		if err != nil {
			fmt.Println("Ошибка подключения к SSH:", err)
			fmt.Println("Попробуйте ввести данные снова.")
			continue
		}
		break
	}

	var snmpService *snmp.SnmpService
	for {
		snmpIP, snmpPort, snmpCommunity := services.GetSNMPInput() // Используем функцию из input_service
		snmpService = snmp.NewSnmpService(snmpIP, snmpPort, snmpCommunity)
		err := snmpService.Connect()
		if err != nil {
			fmt.Println("Ошибка подключения к SNMP:", err)
			fmt.Println("Попробуйте ввести данные снова.")
			continue
		}
		break
	}
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

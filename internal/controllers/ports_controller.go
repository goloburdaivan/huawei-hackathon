package controllers

import (
	"Hackathon/internal/services"
	"Hackathon/internal/views"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
)

type PortController struct {
	pollingService *services.PollingService
	stopChannel    chan bool
}

func NewPortController(pollingService *services.PollingService) *PortController {
	return &PortController{
		pollingService: pollingService,
		stopChannel:    make(chan bool),
	}
}

func (pc *PortController) ShowPortStats() {
	go func() {
		fmt.Println("Нажмите Enter, чтобы вернуться в меню.")
		fmt.Scanln()
		pc.stopChannel <- true
	}()

	for {
		select {
		case <-pc.stopChannel:
			fmt.Println("Возвращаемся в меню...")
			return
		default:
			portStats := pc.pollingService.GetPortStats()
			views.DisplayPortStats(portStats)
		}
	}
}

func (pc *PortController) ExportPortStatsToCSV() {
	file, err := os.Create("result.csv")
	if err != nil {
		log.Fatal("Не удалось создать файл", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	portStats := pc.pollingService.GetPortStats()

	err = writer.Write([]string{
		"Index", "Name", "OID", "InOctets", "OutOctets",
		"InErrors", "OutErrors", "InUcastPkts", "OutUcastPkts",
		"InMulticastPkts", "OutMulticastPkts", "InBroadcastPkts",
		"OutBroadcastPkts", "AdminStatus", "OperStatus",
	})
	if err != nil {
		log.Fatal("Не удалось записать заголовки в файл", err)
	}

	// Записываем данные статистики в CSV
	for _, stat := range portStats {
		record := []string{
			strconv.Itoa(stat.Index),
			stat.Name,
			stat.OID,
			strconv.FormatUint(uint64(stat.InOctets), 10),
			strconv.FormatUint(uint64(stat.OutOctets), 10),
			strconv.FormatUint(uint64(stat.InErrors), 10),
			strconv.FormatUint(uint64(stat.OutErrors), 10),
			strconv.FormatUint(uint64(stat.InUcastPkts), 10),
			strconv.FormatUint(uint64(stat.OutUcastPkts), 10),
			strconv.FormatUint(uint64(stat.InMulticastPkts), 10),
			strconv.FormatUint(uint64(stat.OutMulticastPkts), 10),
			strconv.FormatUint(uint64(stat.InBroadcastPkts), 10),
			strconv.FormatUint(uint64(stat.OutBroadcastPkts), 10),
			stat.AdminStatus,
			stat.OperStatus,
		}
		err := writer.Write(record)
		if err != nil {
			log.Fatal("Не удалось записать данные в файл", err)
		}
	}

	fmt.Println("Данные успешно экспортированы в result.csv")
}

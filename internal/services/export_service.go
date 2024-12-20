package services

import (
	"Hackathon/internal/core/structs"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
)

type ExportService struct{}

func NewExportService() *ExportService {
	return &ExportService{}
}

func (es *ExportService) ExportPortStatsToCSV(portStats []structs.PortInfo) {
	currentTime := time.Now().Format("20060102_150405")
	fileName := fmt.Sprintf("exportdata_%s.csv", currentTime)
	file, err := os.Create(fileName)

	if err != nil {
		fmt.Println("Не удалось создать файл:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{
		"Index", "Name", "OID", "InOctets", "OutOctets",
		"InErrors", "OutErrors", "InUcastPkts", "OutUcastPkts",
		"InMulticastPkts", "OutMulticastPkts", "InBroadcastPkts",
		"OutBroadcastPkts", "AdminStatus", "OperStatus",
	})
	if err != nil {
		fmt.Println("Не удалось записать заголовки в файл:", err)
		return
	}

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
			fmt.Println("Не удалось записать данные в файл:", err)
			return
		}
	}
	fmt.Println("Данные успешно экспортированы в \n", fileName)
}

func (es *ExportService) ExportPortStatsByIndex(portStats []structs.PortInfo, portIndex int) {
	if portIndex < 0 || portIndex >= len(portStats) {
		fmt.Println("Порт с индексом", portIndex, "не найден")
		return
	}

	portInfo := portStats[portIndex]
	fmt.Println("Экспортируем данные для порта с индексом:", portIndex)
	es.ExportPortStatsToCSV([]structs.PortInfo{portInfo})
}

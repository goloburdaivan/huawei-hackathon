package services

import (
	"Hackathon/internal/core/structs"
	"encoding/csv"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestExportPortStatsToCSV(t *testing.T) {
	es := NewExportService()

	portStats := []structs.PortInfo{
		{
			Index:            1,
			Name:             "Port1",
			OID:              "1.1.1",
			InOctets:         1000,
			OutOctets:        2000,
			InErrors:         0,
			OutErrors:        1,
			InUcastPkts:      500,
			OutUcastPkts:     600,
			InMulticastPkts:  100,
			OutMulticastPkts: 150,
			InBroadcastPkts:  200,
			OutBroadcastPkts: 250,
			AdminStatus:      "up",
			OperStatus:       "up",
		},
	}

	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	es.ExportPortStatsToCSV(portStats)

	w.Close()
	os.Stdout = originalStdout
	outputBytes, _ := ioutil.ReadAll(r)
	output := string(outputBytes)

	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) == 0 {
		t.Fatalf("Не удалось получить вывод функции")
	}
	fileName := strings.TrimSpace(lines[len(lines)-1])

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		t.Fatalf("Файл %s не существует", fileName)
	}

	file, err := os.Open(fileName)
	if err != nil {
		t.Fatalf("Не удалось открыть файл %s: %v", fileName, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("Не удалось прочитать CSV-файл: %v", err)
	}

	if len(records) != 2 {
		t.Errorf("Ожидалось 2 строки (заголовок + данные), получили %d", len(records))
	}

	dataRow := records[1]
	expectedValues := []string{
		strconv.Itoa(portStats[0].Index),
		portStats[0].Name,
		portStats[0].OID,
		strconv.FormatUint(uint64(portStats[0].InOctets), 10),
		strconv.FormatUint(uint64(portStats[0].OutOctets), 10),
		strconv.FormatUint(uint64(portStats[0].InErrors), 10),
		strconv.FormatUint(uint64(portStats[0].OutErrors), 10),
		strconv.FormatUint(uint64(portStats[0].InUcastPkts), 10),
		strconv.FormatUint(uint64(portStats[0].OutUcastPkts), 10),
		strconv.FormatUint(uint64(portStats[0].InMulticastPkts), 10),
		strconv.FormatUint(uint64(portStats[0].OutMulticastPkts), 10),
		strconv.FormatUint(uint64(portStats[0].InBroadcastPkts), 10),
		strconv.FormatUint(uint64(portStats[0].OutBroadcastPkts), 10),
		portStats[0].AdminStatus,
		portStats[0].OperStatus,
	}

	if len(dataRow) != len(expectedValues) {
		t.Fatalf("Ожидалось %d столбцов, получили %d", len(expectedValues), len(dataRow))
	}

	for i, expected := range expectedValues {
		if dataRow[i] != expected {
			t.Errorf("В столбце %d ожидалось '%s', получили '%s'", i, expected, dataRow[i])
		}
	}

	os.Remove(fileName)
}

func TestExportPortStatsByIndex(t *testing.T) {
	es := NewExportService()

	portStats := []structs.PortInfo{
		{Index: 0, Name: "Port0"},
		{Index: 1, Name: "Port1"},
	}

	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	es.ExportPortStatsByIndex(portStats, 1)

	w.Close()
	os.Stdout = originalStdout
	outputBytes, _ := ioutil.ReadAll(r)
	output := string(outputBytes)

	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) == 0 {
		t.Fatalf("Не удалось получить вывод функции")
	}
	fileName := strings.TrimSpace(lines[len(lines)-1])

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		t.Fatalf("Файл %s не существует", fileName)
	}

	file, err := os.Open(fileName)
	if err != nil {
		t.Fatalf("Не удалось открыть файл %s: %v", fileName, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("Не удалось прочитать CSV-файл: %v", err)
	}

	if len(records) != 2 {
		t.Errorf("Ожидалось 2 строки (заголовок + данные), получили %d", len(records))
	}

	dataRow := records[1]
	expectedValues := []string{
		strconv.Itoa(portStats[1].Index),
		portStats[1].Name,
		portStats[1].OID,
		"0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "", "",
	}

	if len(dataRow) != len(expectedValues) {
		t.Fatalf("Ожидалось %d столбцов, получили %d", len(expectedValues), len(dataRow))
	}

	for i, expected := range expectedValues {
		if dataRow[i] != expected {
			t.Errorf("В столбце %d ожидалось '%s', получили '%s'", i, expected, dataRow[i])
		}
	}

	os.Remove(fileName)

	r, w, _ = os.Pipe()
	os.Stdout = w

	es.ExportPortStatsByIndex(portStats, 5)

	w.Close()
	os.Stdout = originalStdout
	outputBytes, _ = ioutil.ReadAll(r)
	output = string(outputBytes)

	if !strings.Contains(output, "Порт с индексом 5 не найден") {
		t.Errorf("Ожидалось сообщение об ошибке для неверного индекса, получили '%s'", output)
	}
}
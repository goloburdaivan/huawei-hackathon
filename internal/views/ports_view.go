package views

import (
	"Hackathon/internal/core/structs"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"os"
	"os/exec"
	"runtime"
	"time"
)

func DisplayPortStats(portStats []structs.PortInfo) {
	clearConsole()

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"Index", "Name", "OID", "InOctets", "OutOctets", "InErrors", "OutErrors",
		"InUcastPkts", "OutUcastPkts", "InMulticastPkts", "OutMulticastPkts",
		"InBroadcastPkts", "OutBroadcastPkts", "AdminStatus", "OperStatus",
	})

	for _, port := range portStats {
		row := []string{
			fmt.Sprintf("%d", port.Index),
			port.Name,
			port.OID,
			fmt.Sprintf("%d", port.InOctets),
			fmt.Sprintf("%d", port.OutOctets),
			fmt.Sprintf("%d", port.InErrors),
			fmt.Sprintf("%d", port.OutErrors),
			fmt.Sprintf("%d", port.InUcastPkts),
			fmt.Sprintf("%d", port.OutUcastPkts),
			fmt.Sprintf("%d", port.InMulticastPkts),
			fmt.Sprintf("%d", port.OutMulticastPkts),
			fmt.Sprintf("%d", port.InBroadcastPkts),
			fmt.Sprintf("%d", port.OutBroadcastPkts),
			port.AdminStatus,
			port.OperStatus,
		}
		table.Append(row)
	}

	table.Render()
	fmt.Println("\nОбновление через 1 секунду. Нажмите Enter, чтобы вернуться в меню.")
	time.Sleep(500 * time.Millisecond)
}

func DisplaySinglePortStats(portStat *structs.PortInfo) {
	fmt.Printf("%-20s: %v\n", "Id", portStat.Index)
	fmt.Printf("%-20s: %v\n", "Name", portStat.Name)
	fmt.Printf("%-20s: %v\n", "OID", portStat.OID)
	fmt.Printf("%-20s: %v\n", "InOctets", portStat.InOctets)
	fmt.Printf("%-20s: %v\n", "OutOctets", portStat.OutOctets)
	fmt.Printf("%-20s: %v\n", "InErrors", portStat.InErrors)
	fmt.Printf("%-20s: %v\n", "OutErrors", portStat.OutErrors)
	fmt.Printf("%-20s: %v\n", "InUcastPkts", portStat.InUcastPkts)
	fmt.Printf("%-20s: %v\n", "OutUcastPkts", portStat.OutUcastPkts)
	fmt.Printf("%-20s: %v\n", "InMulticastPkts", portStat.InMulticastPkts)
	fmt.Printf("%-20s: %v\n", "OutMulticastPkts", portStat.OutMulticastPkts)
	fmt.Printf("%-20s: %v\n", "InBroadcastPkts", portStat.InBroadcastPkts)
	fmt.Printf("%-20s: %v\n", "OutBroadcastPkts", portStat.OutBroadcastPkts)
	fmt.Printf("%-20s: %v\n", "InOctetsPkts", portStat.InOctetsPkts)
	fmt.Printf("%-20s: %v\n", "OutOctetsPkts", portStat.OutOctetsPkts)
	fmt.Printf("%-20s: %v\n", "InBandwidthUtil", portStat.InBandwidthUtil)
	fmt.Printf("%-20s: %v\n", "OutBandwidthUtil", portStat.OutBandwidthUtil)
}

func clearConsole() {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	case "linux", "darwin":
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		fmt.Println("Не удалось очистить экран, неизвестная ОС.")
	}
}

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
	fmt.Printf("%-10s %-20s %-10s %-10s %-10s %-10s %-10s %-10s %-10s %-10s %-10s %-10s %-10s %-10s %-10s %-10s\n",
		"Index",
		"Name",
		"OID",
		"InOctets",
		"OutOctets",
		"InErrors",
		"OutErrors",
		"InUcastPkts",
		"OutUcastPkts",
		"InMulticast",
		"OutMulticast",
		"InBroadcast",
		"OutBroadcast",
		"InOctetsPkts",
		"OutOctetsPkts",
		"InBWUtil",
		"OutBWUtil")

	fmt.Printf("%-10d %-20s %-10s %-10d %-10d %-10d %-10d %-10d %-10d %-10d %-10d %-10d %-10d %-10d %-10d %-10.2f %-10.2f\n",
		portStat.Index,
		portStat.Name,
		portStat.OID,
		portStat.InOctets,
		portStat.OutOctets,
		portStat.InErrors,
		portStat.OutErrors,
		portStat.InUcastPkts,
		portStat.OutUcastPkts,
		portStat.InMulticastPkts,
		portStat.OutMulticastPkts,
		portStat.InBroadcastPkts,
		portStat.OutBroadcastPkts,
		portStat.InOctetsPkts,
		portStat.OutOctetsPkts,
		portStat.InBandwidthUtil,
		portStat.OutBandwidthUtil)
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

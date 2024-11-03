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

func DisplayPortList(portStats []structs.PortInfo) {
	clearConsole()

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"Index", "Name",
	})

	for _, port := range portStats {
		row := []string{
			fmt.Sprintf("%d", port.Index),
			port.Name,
		}
		table.Append(row)
	}
	table.Render()

}

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
	fmt.Println("\nUpdating in 1 second. Press Enter to return to the menu.")
	time.Sleep(500 * time.Millisecond)
}

func DisplaySinglePortStats(portStat *structs.PortInfo) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Property", "Value"})
	table.Append([]string{"Id", fmt.Sprintf("%v", portStat.Index)})
	table.Append([]string{"Name", portStat.Name})
	table.Append([]string{"OID", portStat.OID})
	table.Append([]string{"InOctets", fmt.Sprintf("%v", portStat.InOctets)})
	table.Append([]string{"OutOctets", fmt.Sprintf("%v", portStat.OutOctets)})
	table.Append([]string{"InErrors", fmt.Sprintf("%v", portStat.InErrors)})
	table.Append([]string{"OutErrors", fmt.Sprintf("%v", portStat.OutErrors)})
	table.Append([]string{"InUcastPkts", fmt.Sprintf("%v", portStat.InUcastPkts)})
	table.Append([]string{"OutUcastPkts", fmt.Sprintf("%v", portStat.OutUcastPkts)})
	table.Append([]string{"InMulticastPkts", fmt.Sprintf("%v", portStat.InMulticastPkts)})
	table.Append([]string{"OutMulticastPkts", fmt.Sprintf("%v", portStat.OutMulticastPkts)})
	table.Append([]string{"InBroadcastPkts", fmt.Sprintf("%v", portStat.InBroadcastPkts)})
	table.Append([]string{"OutBroadcastPkts", fmt.Sprintf("%v", portStat.OutBroadcastPkts)})
	table.Append([]string{"InOctetsPkts", fmt.Sprintf("%v", portStat.InOctetsPkts)})
	table.Append([]string{"OutOctetsPkts", fmt.Sprintf("%v", portStat.OutOctetsPkts)})
	table.Append([]string{"InBandwidthUtil", fmt.Sprintf("%v", portStat.InBandwidthUtil)})
	table.Append([]string{"OutBandwidthUtil", fmt.Sprintf("%v", portStat.OutBandwidthUtil)})
	table.Append([]string{"InBandwidthActual", fmt.Sprintf("%v", portStat.InBandwidthActual)})
	table.Append([]string{"OutBandwidthActual", fmt.Sprintf("%v", portStat.OutBandwidthActual)})
	table.Render()
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
		fmt.Println("Failed to clear the screen, unknown OS.")
	}
}

package ssh

import (
	"Hackathon/internal/core/structs"
	"strconv"
	"strings"
)

func parseFieldValue(line, delimiter string) string {
	parts := strings.Split(line, delimiter)
	if len(parts) == 2 {
		return strings.TrimSpace(parts[1])
	}
	return ""
}

func parseTrafficData(line string, bytesField, packetsField *uint) {
	parts := strings.Split(line, ":")
	if len(parts) == 2 {
		data := strings.Split(parts[1], ",")
		if len(data) >= 2 {
			*bytesField = parseUint(strings.TrimSuffix(strings.TrimSpace(data[0]), " bytes"))
			*packetsField = parseUint(strings.TrimSuffix(strings.TrimSpace(data[1]), " packets"))
		}
	}
}

func parsePacketCounts(line, section string, currentPort *structs.PortInfo) {
	for _, part := range strings.Split(line, ",") {
		part = strings.TrimSpace(part)
		packetType := strings.Split(part, ":")
		if len(packetType) == 2 {
			pktCount := parseUint(strings.TrimSuffix(strings.TrimSpace(packetType[1]), " packets"))
			switch section {
			case "input":
				setInputCounts(packetType[0], pktCount, currentPort)
			case "output":
				setOutputCounts(packetType[0], pktCount, currentPort)
			}
		}
	}
}

func setInputCounts(packetType string, count uint, currentPort *structs.PortInfo) {
	switch {
	case strings.HasPrefix(packetType, "Unicast"):
		currentPort.InUcastPkts = count
	case strings.HasPrefix(packetType, "Multicast"):
		currentPort.InMulticastPkts = count
	case strings.HasPrefix(packetType, "Broadcast"):
		currentPort.InBroadcastPkts = count
	}
}

func setOutputCounts(packetType string, count uint, currentPort *structs.PortInfo) {
	switch {
	case strings.HasPrefix(packetType, "Unicast"):
		currentPort.OutUcastPkts = count
	case strings.HasPrefix(packetType, "Multicast"):
		currentPort.OutMulticastPkts = count
	case strings.HasPrefix(packetType, "Broadcast"):
		currentPort.OutBroadcastPkts = count
	}
}

func parseFloat(value string) float32 {
	if n, err := strconv.ParseFloat(value, 32); err == nil {
		return float32(n)
	}
	return 0
}

func parseUint(value string) uint {
	if n, err := strconv.Atoi(value); err == nil {
		return uint(n)
	}
	return 0
}

func parseFloatSuffix(line, delimiter string) float32 {
	valueStr := parseFieldValue(line, delimiter)
	return parseFloat(strings.TrimSuffix(valueStr, "%"))
}

func isTrafficLine(line string) bool {
	return !strings.Contains(line, "Unicast") && !strings.Contains(line, "Multicast") && !strings.Contains(line, "Broadcast")
}

func parseState(line string) string {
	parts := strings.Split(line, ":")
	if len(parts) > 1 {
		return strings.TrimSpace(parts[1])
	}
	return ""
}

func hasAnsiEscapeCodes(s string) bool {
	return strings.ContainsRune(s, '\x1b')
}

package ssh

import (
	"Hackathon/internal/core/structs"
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func (s *SshService) parsePorts() {
	s.portStats = nil
	data, _ := s.RunCommand("dis int brief")
	lines := strings.Split(data, "\n")

	startIndex := -1
	for i, line := range lines {
		if strings.Contains(line, "Interface") {
			startIndex = i + 1
			break
		}
	}

	for i := startIndex; i < len(lines); i++ {
		line := strings.Fields(lines[i])

		if len(line) < 7 {
			continue
		}

		name := line[0]

		inUti, _ := strconv.Atoi(strings.TrimSuffix(line[3], "%"))
		outUti, _ := strconv.Atoi(strings.TrimSuffix(line[4], "%"))
		inErrors, _ := strconv.Atoi(line[5])
		outErrors, _ := strconv.Atoi(line[6])

		if !hasAnsiEscapeCodes(name) {
			port := structs.PortInfo{
				Index:        i - startIndex,
				Name:         strings.ReplaceAll(strings.TrimSpace(name), "\n", ""),
				OID:          "Not available for SSH",
				InErrors:     uint(inErrors),
				OutErrors:    uint(outErrors),
				InUcastPkts:  uint(inUti),
				OutUcastPkts: uint(outUti),
				AdminStatus:  "UNKNOWN",
				OperStatus:   "UNKNOWN",
			}
			s.portStats = append(s.portStats, port)
		}
	}

	s.portStats = s.portStats[1 : len(s.portStats)-1]
}

func (s *SshService) parseInterfaceData(port *structs.PortInfo) {
	output, _ := s.RunCommand(fmt.Sprintf("display interface %s", port.Name))
	interfaceRegex := regexp.MustCompile(`^(\S+)\s+current state\s+:\s+(\S+)`)
	scanner := bufio.NewScanner(strings.NewReader(output))
	var currentSection string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		line = strings.ReplaceAll(line, "Administratively", "")
		line = strings.ReplaceAll(line, "(spoofing)", "")

		switch {
		case strings.HasPrefix(line, "Input :") && isTrafficLine(line):
			parseTrafficData(line, &port.InOctets, &port.InOctetsPkts)

		case strings.HasPrefix(line, "Output:") && isTrafficLine(line):
			parseTrafficData(line, &port.OutOctets, &port.OutOctetsPkts)

		case strings.HasPrefix(line, "Input:"):
			currentSection = "input"

		case strings.HasPrefix(line, "Output:"):
			currentSection = "output"

		case strings.HasPrefix(line, "Unicast:") || strings.HasPrefix(line, "Multicast:") || strings.HasPrefix(line, "Broadcast:"):
			parsePacketCounts(line, currentSection, port)

		case strings.HasPrefix(line, "Input bandwidth utilization threshold"):
			port.InBandwidthUtil = parseFloatSuffix(line, " :")

		case strings.HasPrefix(line, "Output bandwidth utilization threshold"):
			port.OutBandwidthUtil = parseFloatSuffix(line, ":")
		case interfaceRegex.MatchString(line):
			matches := interfaceRegex.FindStringSubmatch(line)
			port.OperStatus = matches[2]
		case strings.HasPrefix(line, "Last 300 seconds input utility rate"):
			port.InBandwidthActual = parseFloatSuffix(line, ":")

		case strings.HasPrefix(line, "Last 300 seconds output utility rate"):
			port.OutBandwidthActual = parseFloatSuffix(line, ":")

		case strings.Contains(line, "Line protocol current state"):
			port.AdminStatus = parseState(line)
		}
	}
}

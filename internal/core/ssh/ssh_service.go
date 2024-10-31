package ssh

import (
	"Hackathon/internal/core/structs"
	"Hackathon/internal/services"
	"bufio"
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

type SshService struct {
	Target     string
	Port       uint16
	User       string
	Password   string
	connection *ssh.Client
	session    *ssh.Session
	stdinPipe  io.WriteCloser
	stdoutPipe io.Reader
	portStats  []structs.PortInfo
	mu         sync.Mutex
}

func (s *SshService) GetPortStats() []structs.PortInfo {
	return s.portStats
}

func (s *SshService) PollStatistics() error {
	s.parsePorts()

	for i := range s.portStats {
		s.parseInterfaceData(&s.portStats[i])
	}

	return nil
}

func (s *SshService) Connect() error {
	config := &ssh.ClientConfig{
		User: s.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.Password),
		},
		Config: ssh.Config{
			KeyExchanges: []string{
				"diffie-hellman-group1-sha1",
				"diffie-hellman-group-exchange-sha1",
			},
			Ciphers: []string{
				"aes128-cbc",
				"3des-cbc",
				"des-cbc",
			},
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Duration(2) * time.Second,
	}

	connection, err := s.initConnection(config)
	if err != nil {
		return err
	}

	s.connection = connection
	s.session, err = s.connection.NewSession()

	if err != nil {
		return err
	}

	s.stdinPipe, err = s.session.StdinPipe()

	if err != nil {
		return err
	}

	s.stdoutPipe, err = s.session.StdoutPipe()

	if err != nil {
		return err
	}

	err = s.session.Shell()
	if err != nil {
		return err
	}

	return nil
}

func ConnectSSH() *SshService {
	var sshService *SshService
	for {
		sshIP, sshPort, sshUser, sshPass := services.GetSSHInput()
		sshService = NewSshService(sshIP, sshPort, sshUser, sshPass)
		err := sshService.Connect()
		if err != nil {
			fmt.Println("Ошибка подключения к SSH:", err)
			fmt.Println("Попробуйте ввести данные снова.")
			continue
		}
		break
	}
	return sshService
}

func (s *SshService) initConnection(config *ssh.ClientConfig) (*ssh.Client, error) {
	return ssh.Dial("tcp", fmt.Sprintf("%s:%d", s.Target, s.Port), config)
}

func (s *SshService) CloseConnection() {
	if s.connection != nil {
		fmt.Println("Closing SSH connection")
		s.connection.Close()
		s.session.Close()
		s.stdinPipe.Close()
	}
}

func NewSshService(target string, port uint16, user, password string) *SshService {
	return &SshService{
		Target:   target,
		Port:     port,
		User:     user,
		Password: password,
	}
}

func (s *SshService) RunCommand(cmd string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	commands := []string{
		"screen-length 0 temporary",
		"sys",
		cmd,
		"echo END_COMMAND",
	}

	collectOutput := false
	var outputBuf bytes.Buffer

	for _, command := range commands {
		_, err := fmt.Fprintf(s.stdinPipe, "%s\n", command)
		if err != nil {
			return "", fmt.Errorf("failed to send command: %w", err)
		}
		if command == cmd {
			collectOutput = true
		}
		time.Sleep(500 * time.Millisecond)
	}

	buf := make([]byte, 1024)
	for {
		n, err := s.stdoutPipe.Read(buf)
		if n > 0 {
			output := string(buf[:n])
			if collectOutput {
				outputBuf.Write(buf[:n])
				if strings.Contains(output, "END_COMMAND") {
					break
				}
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", fmt.Errorf("failed to read stdout: %w", err)
		}
	}

	result := outputBuf.String()
	result = strings.Replace(result, "END_COMMAND", "", -1)
	result = strings.TrimSpace(result)

	return result, nil
}

func (s *SshService) parsePorts() {
	s.portStats = nil
	data, _ := s.RunCommand("dis int brief")
	lines := strings.Split(data, "\n")

	for i := 8; i < len(lines); i++ {
		line := strings.Fields(lines[i])

		if len(line) < 6 {
			continue
		}

		name := line[0]
		adminStatus := line[1]
		operStatus := line[2]
		inUti, _ := strconv.Atoi(strings.TrimSuffix(line[3], "%"))
		outUti, _ := strconv.Atoi(strings.TrimSuffix(line[4], "%"))
		inErrors, _ := strconv.Atoi(line[5])
		outErrors, _ := strconv.Atoi(line[6])

		port := structs.PortInfo{
			Index:        i - 7,
			Name:         name,
			OID:          "Not available for SSH",
			InErrors:     uint(inErrors),
			OutErrors:    uint(outErrors),
			InUcastPkts:  uint(inUti),
			OutUcastPkts: uint(outUti),
			AdminStatus:  adminStatus,
			OperStatus:   operStatus,
		}

		s.portStats = append(s.portStats, port)
	}

	s.portStats = s.portStats[1 : len(s.portStats)-1]
}

func (s *SshService) parseInterfaceData(port *structs.PortInfo) {
	output, _ := s.RunCommand(fmt.Sprintf("display interface %s", port.Name))
	scanner := bufio.NewScanner(strings.NewReader(output))
	var currentSection string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		switch {
		case strings.HasPrefix(line, "Input:") && isTrafficLine(line):
			parseTrafficData(line, &port.InOctets, &port.InOctetsPkts)

		case strings.HasPrefix(line, "Output:") && isTrafficLine(line):
			parseTrafficData(line, &port.OutOctets, &port.OutOctetsPkts)

		case strings.HasPrefix(line, "Input:"):
			currentSection = "input"

		case strings.HasPrefix(line, "Output:"):
			currentSection = "output"

		case strings.HasPrefix(line, "Unicast:") || strings.HasPrefix(line, "Multicast:") || strings.HasPrefix(line, "Broadcast:"):
			parsePacketCounts(line, currentSection, port)

		case strings.HasPrefix(line, "Input bandwidth utilization"):
			port.InBandwidthUtil = parseUintSuffix(line, ":")

		case strings.HasPrefix(line, "Output bandwidth utilization"):
			port.OutBandwidthUtil = parseUintSuffix(line, ":")
		}
	}
}

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

func (s *SshService) ParseDeviceStatus() (structs.DeviceStatus, error) {
	data, err := s.RunCommand("display device")
	lines := strings.Split(data, "\n")

	if err != nil {
		log.Fatalf("ошибка при выполнении команды: %v", err)
	}

	if len(lines) < 3 {
		return structs.DeviceStatus{}, fmt.Errorf("not enough data to parse")
	}

	fields := strings.Fields(lines[3])

	if len(fields) < 8 {
		return structs.DeviceStatus{}, fmt.Errorf("unexpected number of fields")
	}

	slot, err := strconv.Atoi(fields[0])
	if err != nil {
		return structs.DeviceStatus{}, fmt.Errorf("error parsing slot number: %v", err)
	}

	return structs.DeviceStatus{
		Slot:     slot,
		Sub:      fields[1],
		Type:     fields[2],
		Online:   fields[3],
		Power:    fields[4],
		Register: fields[5],
		Status:   fields[6],
		Role:     fields[7],
	}, nil
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

func parseUint(value string) uint {
	if n, err := strconv.Atoi(value); err == nil {
		return uint(n)
	}
	return 0
}

func parseUintSuffix(line, delimiter string) uint {
	valueStr := parseFieldValue(line, delimiter)
	return parseUint(strings.TrimSuffix(valueStr, "%"))
}

func isTrafficLine(line string) bool {
	return !strings.Contains(line, "Unicast") && !strings.Contains(line, "Multicast") && !strings.Contains(line, "Broadcast")
}

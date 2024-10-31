package snmp

import (
	"Hackathon/internal/core/structs"
	"Hackathon/internal/services"
	"fmt"
	"github.com/gosnmp/gosnmp"
	"time"
)

type SnmpService struct {
	Target     string
	Port       uint16
	Community  string
	connection *gosnmp.GoSNMP
	portStats  []structs.PortInfo
}

func (s *SnmpService) PollStatistics() error {
	s.pollPortsStatuses()

	for i := range s.portStats {
		inOctets, _ := getInOctets(s.connection, s.portStats[i].Index)
		outOctets, _ := getOutOctets(s.connection, s.portStats[i].Index)
		inErrors, _ := getInErrors(s.connection, s.portStats[i].Index)
		outErrors, _ := getOutErrors(s.connection, s.portStats[i].Index)
		inUnicast, _ := getInUnicastPackets(s.connection, s.portStats[i].Index)
		outUnicast, _ := getOutUnicastPackets(s.connection, s.portStats[i].Index)
		inMulticast, _ := getInMulticastPackets(s.connection, s.portStats[i].Index)
		outMulticast, _ := getOutMulticastPackets(s.connection, s.portStats[i].Index)
		inBroadcast, _ := getInBroadcastPackets(s.connection, s.portStats[i].Index)
		outBroadcast, _ := getOutBroadcastPackets(s.connection, s.portStats[i].Index)

		s.portStats[i].InOctets = inOctets
		s.portStats[i].OutOctets = outOctets
		s.portStats[i].InErrors = inErrors
		s.portStats[i].OutErrors = outErrors
		s.portStats[i].InUcastPkts = inUnicast
		s.portStats[i].OutUcastPkts = outUnicast
		s.portStats[i].InMulticastPkts = inMulticast
		s.portStats[i].OutMulticastPkts = outMulticast
		s.portStats[i].InBroadcastPkts = inBroadcast
		s.portStats[i].OutBroadcastPkts = outBroadcast
	}

	return nil
}

func (s *SnmpService) GetPortStats() []structs.PortInfo {
	return s.portStats
}

func (s *SnmpService) Connect() error {
	return s.connection.Connect()
}

func (s *SnmpService) FetchPorts() error {
	s.portStats = []structs.PortInfo{}
	indexes, err := getPortsIndexes(s.connection)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	descriptions, err := getPortsDescriptions(s.connection)
	if err != nil {
		return nil
	}

	s.portStats = initPorts(indexes, descriptions)

	return nil
}

func ConnectSNMP() (*SnmpService, error) {
	var snmpService *SnmpService
	for {
		snmpIP, snmpPort, snmpCommunity := services.GetSNMPInput()
		snmpService = NewSnmpService(snmpIP, snmpPort, snmpCommunity)
		err := snmpService.Connect()
		if err != nil {
			return nil, fmt.Errorf("Ошибка подключения к SNMP: %v", err)
		}
		break
	}
	return snmpService, nil
}

func (s *SnmpService) CloseConnection() {
	if s.connection != nil {
		fmt.Println("Closing SNMP connection")
		s.connection.Conn.Close()
	}
}

func NewSnmpService(target string, port uint16, community string) *SnmpService {
	return &SnmpService{
		Target:    target,
		Port:      port,
		Community: community,
		connection: &gosnmp.GoSNMP{
			Target:    target,
			Port:      port,
			Community: community,
			Version:   gosnmp.Version2c,
			Timeout:   time.Duration(2) * time.Second,
			Retries:   1,
		},
	}
}

func (s *SnmpService) pollPortsStatuses() {
	adminStatuses, _ := getAdminStatuses(s.connection)
	operStatuses, _ := getOperStatuses(s.connection)

	for i := range s.portStats {
		s.portStats[i].AdminStatus = getStatusLabel(adminStatuses[i].Value.(int))
		s.portStats[i].OperStatus = getStatusLabel(operStatuses[i].Value.(int))
	}
}

func initPorts(indexes []gosnmp.SnmpPDU, descriptions []gosnmp.SnmpPDU) []structs.PortInfo {
	ports := make([]structs.PortInfo, len(indexes))
	for i := range indexes {
		ports[i] = structs.PortInfo{
			Index: indexes[i].Value.(int),
			Name:  string(descriptions[i].Value.([]byte)),
			OID:   indexes[i].Name,
		}
	}

	return ports
}

package snmp

import (
	"Hackathon/internal/core/structs"
	"fmt"
	"github.com/gosnmp/gosnmp"
	"time"
)

type SnmpService struct {
	Target     string
	Port       uint16
	Community  string
	connection *gosnmp.GoSNMP
	PortStats  []structs.PortInfo
}

func (s *SnmpService) PollStatistics() error {
	s.pollPortsStatuses()

	for i := range s.PortStats {
		inOctets, _ := getInOctets(s.connection, s.PortStats[i].Index)
		outOctets, _ := getOutOctets(s.connection, s.PortStats[i].Index)
		inErrors, _ := getInErrors(s.connection, s.PortStats[i].Index)
		outErrors, _ := getOutErrors(s.connection, s.PortStats[i].Index)
		inUnicast, _ := getInUnicastPackets(s.connection, s.PortStats[i].Index)
		outUnicast, _ := getOutUnicastPackets(s.connection, s.PortStats[i].Index)
		inMulticast, _ := getInMulticastPackets(s.connection, s.PortStats[i].Index)
		outMulticast, _ := getOutMulticastPackets(s.connection, s.PortStats[i].Index)
		inBroadcast, _ := getInBroadcastPackets(s.connection, s.PortStats[i].Index)
		outBroadcast, _ := getOutBroadcastPackets(s.connection, s.PortStats[i].Index)

		s.PortStats[i].InOctets = inOctets
		s.PortStats[i].OutOctets = outOctets
		s.PortStats[i].InErrors = inErrors
		s.PortStats[i].OutErrors = outErrors
		s.PortStats[i].InUcastPkts = inUnicast
		s.PortStats[i].OutUcastPkts = outUnicast
		s.PortStats[i].InMulticastPkts = inMulticast
		s.PortStats[i].OutMulticastPkts = outMulticast
		s.PortStats[i].InBroadcastPkts = inBroadcast
		s.PortStats[i].OutBroadcastPkts = outBroadcast
	}

	return nil
}

func (s *SnmpService) Connect() error {
	return s.connection.Connect()
}

func (s *SnmpService) FetchPorts() error {
	s.PortStats = []structs.PortInfo{}
	indexes, err := getPortsIndexes(s.connection)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	descriptions, err := getPortsDescriptions(s.connection)
	if err != nil {
		return nil
	}

	s.PortStats = initPorts(indexes, descriptions)

	return nil
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

	for i := range s.PortStats {
		s.PortStats[i].AdminStatus = getStatusLabel(adminStatuses[i].Value.(int))
		s.PortStats[i].OperStatus = getStatusLabel(operStatuses[i].Value.(int))
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

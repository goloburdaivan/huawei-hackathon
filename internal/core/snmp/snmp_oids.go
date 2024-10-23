package snmp

import (
	"fmt"
	"github.com/gosnmp/gosnmp"
)

func getPortsIndexes(connection *gosnmp.GoSNMP) ([]gosnmp.SnmpPDU, error) {
	return connection.BulkWalkAll("1.3.6.1.2.1.2.2.1.1")
}

func getPortsDescriptions(connection *gosnmp.GoSNMP) ([]gosnmp.SnmpPDU, error) {
	return connection.BulkWalkAll("1.3.6.1.2.1.2.2.1.2")
}

func getAdminStatuses(connection *gosnmp.GoSNMP) ([]gosnmp.SnmpPDU, error) {
	return connection.BulkWalkAll("1.3.6.1.2.1.2.2.1.7")
}

func getOperStatuses(connection *gosnmp.GoSNMP) ([]gosnmp.SnmpPDU, error) {
	return connection.BulkWalkAll("1.3.6.1.2.1.2.2.1.8")
}

func getInOctets(connection *gosnmp.GoSNMP, interfaceIndex int) (uint, error) {
	result, err := connection.Get([]string{fmt.Sprintf("%s.%d", "1.3.6.1.2.1.2.2.1.10", interfaceIndex)})
	return result.Variables[0].Value.(uint), err
}

func getOutOctets(connection *gosnmp.GoSNMP, interfaceIndex int) (uint, error) {
	result, err := connection.Get([]string{fmt.Sprintf("%s.%d", "1.3.6.1.2.1.2.2.1.16", interfaceIndex)})
	return result.Variables[0].Value.(uint), err
}

func getInErrors(connection *gosnmp.GoSNMP, interfaceIndex int) (uint, error) {
	result, err := connection.Get([]string{fmt.Sprintf("%s.%d", "1.3.6.1.2.1.2.2.1.14", interfaceIndex)})

	return result.Variables[0].Value.(uint), err
}

func getOutErrors(connection *gosnmp.GoSNMP, interfaceIndex int) (uint, error) {
	result, err := connection.Get([]string{fmt.Sprintf("%s.%d", "1.3.6.1.2.1.2.2.1.20", interfaceIndex)})

	return result.Variables[0].Value.(uint), err
}

func getInUnicastPackets(connection *gosnmp.GoSNMP, interfaceIndex int) (uint, error) {
	result, err := connection.Get([]string{fmt.Sprintf("%s.%d", "1.3.6.1.2.1.2.2.1.11", interfaceIndex)})

	return result.Variables[0].Value.(uint), err
}

func getOutUnicastPackets(connection *gosnmp.GoSNMP, interfaceIndex int) (uint, error) {
	result, err := connection.Get([]string{fmt.Sprintf("%s.%d", "1.3.6.1.2.1.2.2.1.17", interfaceIndex)})

	return result.Variables[0].Value.(uint), err
}

func getInMulticastPackets(connection *gosnmp.GoSNMP, interfaceIndex int) (uint, error) {
	result, err := connection.Get([]string{fmt.Sprintf("%s.%d", "1.3.6.1.2.1.2.2.1.12", interfaceIndex)})

	return result.Variables[0].Value.(uint), err
}

func getOutMulticastPackets(connection *gosnmp.GoSNMP, interfaceIndex int) (uint, error) {
	result, err := connection.Get([]string{fmt.Sprintf("%s.%d", "1.3.6.1.2.1.2.2.1.18", interfaceIndex)})

	return result.Variables[0].Value.(uint), err
}

func getInBroadcastPackets(connection *gosnmp.GoSNMP, interfaceIndex int) (uint, error) {
	result, err := connection.Get([]string{fmt.Sprintf("%s.%d", "1.3.6.1.2.1.2.2.1.13", interfaceIndex)})

	return result.Variables[0].Value.(uint), err
}

func getOutBroadcastPackets(connection *gosnmp.GoSNMP, interfaceIndex int) (uint, error) {
	result, err := connection.Get([]string{fmt.Sprintf("%s.%d", "1.3.6.1.2.1.2.2.1.19", interfaceIndex)})

	return result.Variables[0].Value.(uint), err
}

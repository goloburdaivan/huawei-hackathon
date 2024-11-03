package services

import (
	"Hackathon/internal/core/snmp"
	"Hackathon/internal/core/ssh"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ConnectSSH() *ssh.SshService {
	sshIP, sshPort, sshUser, sshPass := getSSHCredentials()
	sshService := trySSHConnection(sshIP, sshPort, sshUser, sshPass)
	return sshService
}

func getSSHCredentials() (string, uint16, string, string) {
	var sshIP, sshUser, sshPass string
	var sshPort uint16

	if fileExists("last_login") {
		fmt.Println("Use saved data from last_login? (y/n)")
		var choice string
		fmt.Scan(&choice)
		if strings.ToLower(choice) == "y" {
			sshIP, sshPort, sshUser, sshPass = readLastLogin()
		} else {
			sshIP, sshPort, sshUser, sshPass = getAndSaveSSHInput()
		}
	} else {
		sshIP, sshPort, sshUser, sshPass = getAndSaveSSHInput()
	}

	return sshIP, sshPort, sshUser, sshPass
}

func trySSHConnection(sshIP string, sshPort uint16, sshUser, sshPass string) *ssh.SshService {
	var sshService *ssh.SshService

	for {
		sshService = ssh.NewSshService(sshIP, sshPort, sshUser, sshPass)
		err := sshService.Connect()
		if err != nil {
			fmt.Println("SSH connection error:", err)
			fmt.Println("Please try entering the data again.")
			sshIP, sshPort, sshUser, sshPass = getAndSaveSSHInput()
			continue
		}
		break
	}

	return sshService
}

func getAndSaveSSHInput() (string, uint16, string, string) {
	sshIP, sshPort, sshUser, sshPass := GetSSHInput()
	saveLastLogin(sshIP, sshPort, sshUser, sshPass)
	return sshIP, sshPort, sshUser, sshPass
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func readLastLogin() (string, uint16, string, string) {
	file, err := os.Open("last_login")
	if err != nil {
		fmt.Println("Error reading the last_login file:", err)
		return "", 0, "", ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var data []string
	for scanner.Scan() {
		data = append(data, scanner.Text())
	}

	if len(data) < 4 {
		fmt.Println("Not enough data in last_login, please enter manually.")
		return "", 0, "", ""
	}

	port, err := strconv.ParseUint(data[1], 10, 16)
	if err != nil {
		fmt.Println("Error converting port from last_login:", err)
		return "", 0, "", ""
	}

	return data[0], uint16(port), data[2], data[3]
}

func saveLastLogin(sshIP string, sshPort uint16, sshUser, sshPass string) {
	file, err := os.Create("last_login")
	if err != nil {
		fmt.Println("Error creating last_login file:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	fmt.Fprintln(writer, sshIP)
	fmt.Fprintln(writer, sshPort)
	fmt.Fprintln(writer, sshUser)
	fmt.Fprintln(writer, sshPass)
	writer.Flush()
}

func ConnectSNMP() *snmp.SnmpService {
	var snmpService *snmp.SnmpService
	for {
		snmpIP, snmpPort, snmpCommunity := GetSNMPInput()
		snmpService = snmp.NewSnmpService(snmpIP, snmpPort, snmpCommunity)
		err := snmpService.Connect()
		if err != nil {
			fmt.Println("Error connecting to SNMP:", err)
			fmt.Println("Please try entering the data again.")
			continue
		}
		break
	}
	return snmpService
}

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
		fmt.Println("Использовать сохраненные данные из last_login? (y/n)")
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
			fmt.Println("Ошибка подключения к SSH:", err)
			fmt.Println("Попробуйте ввести данные снова.")
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
		fmt.Println("Ошибка чтения файла last_login:", err)
		return "", 0, "", ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var data []string
	for scanner.Scan() {
		data = append(data, scanner.Text())
	}

	if len(data) < 4 {
		fmt.Println("Недостаточно данных в last_login, введите вручную.")
		return "", 0, "", ""
	}

	port, err := strconv.ParseUint(data[1], 10, 16)
	if err != nil {
		fmt.Println("Ошибка преобразования порта из last_login:", err)
		return "", 0, "", ""
	}

	return data[0], uint16(port), data[2], data[3]
}

func saveLastLogin(sshIP string, sshPort uint16, sshUser, sshPass string) {
	file, err := os.Create("last_login")
	if err != nil {
		fmt.Println("Ошибка создания файла last_login:", err)
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
			fmt.Println("Ошибка подключения к SNMP:", err)
			fmt.Println("Попробуйте ввести данные снова.")
			continue
		}
		break
	}
	return snmpService
}

package services

import (
	"Hackathon/internal/core/snmp"
	"Hackathon/internal/core/ssh"
	"fmt"
)

func ConnectSSH() *ssh.SshService {
	var sshService *ssh.SshService
	for {
		sshIP, sshPort, sshUser, sshPass := GetSSHInput()
		sshService = ssh.NewSshService(sshIP, sshPort, sshUser, sshPass)
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

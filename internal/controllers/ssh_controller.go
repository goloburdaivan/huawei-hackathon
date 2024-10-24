package controllers

import (
	"Hackathon/internal/services"
	"fmt"
)

type SSHController struct {
	service *services.SSHService
}

func NewSSHController(host, port, username, password string) *SSHController {
	sshService := services.NewSSHService(host, port, username, password)
	return &SSHController{
		service: sshService,
	}
}

func (sc *SSHController) StartSSHSession() {
	err := sc.service.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer sc.service.Close()

	fmt.Println("SSH-сессия открыта. Для выхода наберите 'exit'.")

	err = sc.service.StartShell()
	if err != nil {
		fmt.Println(err)
	}
}

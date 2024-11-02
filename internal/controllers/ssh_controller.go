package controllers

import (
	"Hackathon/internal/core/ssh"
)

type SSHController struct {
	sshService *ssh.SshService
}

func NewSSHController(sshService *ssh.SshService) *SSHController {
	return &SSHController{
		sshService: sshService,
	}
}

func (sc *SSHController) StartCliSession() {
	go sc.sshService.StartCliSession()
}

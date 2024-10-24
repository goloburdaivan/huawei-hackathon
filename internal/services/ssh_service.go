package services

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh"
	"os"
	"strings"
)

type SSHService struct {
	Host     string
	Port     string
	Username string
	Password string
	conn     *ssh.Client
	session  *ssh.Session
}

func NewSSHService(host, port, username, password string) *SSHService {
	return &SSHService{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}
}

func (s *SSHService) Connect() error {
	config := &ssh.ClientConfig{
		User: s.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	var err error
	s.conn, err = ssh.Dial("tcp", fmt.Sprintf("%s:%s", s.Host, s.Port), config)
	if err != nil {
		return fmt.Errorf("Ошибка подключения к SSH: %w", err)
	}

	s.session, err = s.conn.NewSession()
	if err != nil {
		return fmt.Errorf("Ошибка создания SSH сессии: %w", err)
	}

	s.session.Stdout = os.Stdout
	s.session.Stderr = os.Stderr
	return nil
}

func (s *SSHService) StartShell() error {
	stdin, err := s.session.StdinPipe()
	if err != nil {
		return fmt.Errorf("Ошибка подключения к STDIN: %w", err)
	}

	err = s.session.Shell()
	if err != nil {
		return fmt.Errorf("Ошибка запуска оболочки: %w", err)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		command, _ := reader.ReadString('\n')
		command = strings.TrimSpace(command)

		if command == "exit" {
			fmt.Println("Завершение SSH-сессии...")
			break
		}

		_, err := stdin.Write([]byte(command + "\n"))
		if err != nil {
			return fmt.Errorf("Ошибка отправки команды: %w", err)
		}
	}

	return nil
}

func (s *SSHService) Close() {
	s.session.Close()
	s.conn.Close()
}

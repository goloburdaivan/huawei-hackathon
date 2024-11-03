package ssh

import (
	"Hackathon/internal/core/structs"
	"bufio"
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
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
	for i := range s.portStats {
		portStat := &s.portStats[i]
		go s.parseInterfaceData(portStat)
	}

	return nil
}

func (s *SshService) FetchPorts() error {
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
				"diffie-hellman-group-exchange-sha256",
			},
			Ciphers: []string{
				"aes256-gcm@openssh.com",
				"aes128-cbc",
				"3des-cbc",
				"aes256-ctr",
			},
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		HostKeyAlgorithms: []string{
			"x509v3-ssh-rsa",
			"ssh-rsa",
			"rsa-sha2-256",
			"rsa-sha2-512",
		},
		Timeout: time.Duration(2) * time.Second,
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

	_, err := fmt.Fprintf(s.stdinPipe, "N\n")
	if err != nil {
		return "", fmt.Errorf("failed to send 'N' to skip password change prompt: %w", err)
	}
	time.Sleep(500 * time.Millisecond)

	var outputBuf bytes.Buffer

	fmt.Fprintf(s.stdinPipe, "%s\n", cmd)

	buf := make([]byte, 1024)
	for {
		n, err := s.stdoutPipe.Read(buf)
		if n > 0 {
			output := string(buf[:n])
			if strings.Contains(output, "---- More ----") {
				_, err = fmt.Fprintf(s.stdinPipe, " ")
				if err != nil {
					return "", fmt.Errorf("failed to send space to continue pagination: %w", err)
				}
				continue
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", fmt.Errorf("failed to read stdout: %w", err)
		}

		outputBuf.Write(buf[:n])

		if n == 16 || n == 18 {
			break
		}
	}

	result := outputBuf.String()

	return result, nil
}

func (s *SshService) StartCliSession() {
	fmt.Println("Starting CLI session. Type 'exit' to quit.")

	for {
		fmt.Print("CLI> ")
		reader := bufio.NewReader(os.Stdin)
		command, _ := reader.ReadString('\n')
		command = strings.TrimSpace(command)

		if command == "exit" {
			fmt.Println("Exiting CLI session.")
			break
		}
		output, err := s.RunCommand(command)
		if err != nil {
			fmt.Printf("Error executing command: %v\n", err)
			continue
		}
		fmt.Println(output)
	}
}

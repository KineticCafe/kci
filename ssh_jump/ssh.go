// Package ssh_jump provides methods for connecting to servers using SSH bastions.
package ssh_jump

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// SSHJump represents an SSH client accessed through a bastion or jump host.
type SSHJump struct {
	JumpHost   string
	JumpUser   string
	TargetHost string
	TargetUser string

	jumpConnection   *ssh.Client
	targetConnection *ssh.Client

	Client *ssh.Client
}

// Helper method for constructing a new SSHJump struct.
func New(jump_host string, jump_user string, target_host string, target_user string) *SSHJump {
	return &SSHJump{
		JumpHost:   jump_host,
		JumpUser:   jump_user,
		TargetHost: target_host,
		TargetUser: target_user,
	}
}

// Connect to the target server through the bastion server. Creates the Client field of the struct.
func (sj *SSHJump) Connect() error {
	socket := os.Getenv("SSH_AUTH_SOCK")
	conn, err := net.Dial("unix", socket)
	if err != nil {
		return fmt.Errorf("Failed to open SSH_AUTH_SOCK: %w", err)
	}

	agentClient := agent.NewClient(conn)
	sshConfig := &ssh.ClientConfig{
		User: sj.JumpUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeysCallback(agentClient.Signers),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	jumpConn, err := ssh.Dial("tcp", sj.JumpHost+":22", sshConfig)
	if err != nil {
		log.Fatalf("Failed to SSH to jump host: %v", err)
	}

	sshConfig.User = sj.TargetUser
	conn, err = jumpConn.Dial("tcp", sj.TargetHost+":22")
	if err != nil {
		return fmt.Errorf("Failed to open target host: %w", err)
	}

	ncc, chans, reqs, err := ssh.NewClientConn(conn, sj.TargetHost+":22", sshConfig)
	if err != nil {
		return fmt.Errorf("Failed to launch client connection: %w", err)
	}

	sj.Client = ssh.NewClient(ncc, chans, reqs)
	return nil
}

// Close all SSH connections if open.
func (sj *SSHJump) Close() {
	if sj.jumpConnection != nil {
		sj.jumpConnection.Close()
	}
	if sj.targetConnection != nil {
		sj.targetConnection.Close()
	}
}

// ExecuteSSHCommands runs a series of commands agains an SSH connection
// and returns the command output in a slice of string.
func ExecuteSSHCommands(conn *ssh.Client, commands []string) ([]string, error) {

	var outputs []string

	for _, command := range commands {
		session, err := conn.NewSession()
		if err != nil {
			return nil, fmt.Errorf("Failed to launch session: %w", err)
		}
		defer session.Close()

		output, err := session.CombinedOutput(command)
		if err != nil {
			return nil, fmt.Errorf("Could not run command %v: %w", command, err)
		}

		outputs = append(outputs, string(output))
	}

	return outputs, nil
}

package ec2_instance

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/KineticCommerce/kci/ssh_jump"
	"golang.org/x/crypto/ssh"
)

// Scan recursively calls Scan on all Instance items. Connection is through a jump
// bastion which needs the connection info, jump and user.
//
// TODO this could easily be handles as a goroutine to speed things up
func (mgr *EC2InstanceManager) JumpScan(host string, user string) error {
	for i := range mgr.Instances {
		err := mgr.Instances[i].JumpScan(host, user)
		if err != nil {
			// We want to continue here...
			log.Printf("target scan failed (%s): %v", mgr.Instances[i].Name, err)
			mgr.Instances[i].Status = "No Connection"
		}
	}

	return nil
}

// JumpScan connects to the underlying instance through a jump host (bastion) and
// runs a Scan.
func (instance *EC2Instance) JumpScan(host string, user string) error {
	jump := ssh_jump.New(host, user, instance.PrivateIP, "ubuntu")

	err := jump.Connect()
	if err != nil {
		return fmt.Errorf("could not perform JumpScan: %w", err)
	}
	defer jump.Close()

	return instance.Scan(jump.Client)
}

// Scan runs a scan over an SSH connection. This will be less important when SSM
// is fully implemented.
func (instance *EC2Instance) Scan(client *ssh.Client) error {
	commands := []string{
		`lsb_release -d | cut -f2 | awk '{print $2}'`,
		`uptime | awk '{print $3 " "  $4}' | tr -d ','`,
		"if [ -f /var/run/reboot-required ]; then echo 'reboot required'; else echo 'no'; fi",
		"sudo cat /var/lib/update-notifier/updates-available | grep 'security updates' | cut -d' ' -f1",
	}

	output, err := ssh_jump.ExecuteSSHCommands(client, commands)
	if err != nil {
		return fmt.Errorf("Error running commands: %w", err)
	}

	instance.OsVersion = strings.ReplaceAll(output[0], "\n", "")
	instance.Uptime = strings.ReplaceAll(output[1], "\n", "")
	instance.RebootRequired = strings.ReplaceAll(output[2], "\n", "")

	updates, err := strconv.Atoi(output[3])
	if err != nil {
		updates = 0 // The error is always an empty string
	}
	instance.SecurityUpdates = updates

	return nil
}

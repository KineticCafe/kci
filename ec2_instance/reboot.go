// Package ec2_instance provides helpers for the AWS SDK v2.
package ec2_instance

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

// Reboot sends a reboot signal to the ec2 instance specified by instanceID. Returns an error
// if the command was unsuccessful. nil, otherwise.
func Reboot(instanceID string) error {
	// TODO extract the config and client load
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return fmt.Errorf("unable to load AWS config, %v", err)
	}

	// TODO inject the client instead of creating here to allow for some better testing
	client := ec2.NewFromConfig(cfg)

	input := &ec2.RebootInstancesInput{
		InstanceIds: []string{
			instanceID,
		},
	}

	_, err = client.RebootInstances(context.Background(), input)

	return err
}

// Package aws_instance provides AWS ec2 helper methods.
package ec2_instance

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// EC2Instance represents an AWS ec2 instance.
type EC2Instance struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	AMI_ID          string `json:"ami_id"`
	InstanceAge     string `json:"instance_age"`
	AMI_Age         string `json:"ami_age"`
	IsSSM           bool   `json:"is_ssm"`
	Status          string `json:"status"`
	PublicIP        string `json:"public_ip"`
	PrivateIP       string `json:"private_ip"`
	OsVersion       string `json:"os_version"`
	RebootRequired  string `json:"reboot_required"`
	SecurityUpdates int    `json:"security_updates"`
	Uptime          string `json:"uptime"`
}

// EC2InstanceManager provides access to a list of EC2 instances. This includes
// methods to fetch, describe, and such.
type EC2InstanceManager struct {
	Instances []EC2Instance
	Client    *ec2.Client
}

// NewManagerWithClient creates a new EC2InstanceManager with a supplied aws client.
func NewManagerWithClient(client *ec2.Client) *EC2InstanceManager {
	return &EC2InstanceManager{
		Instances: []EC2Instance{},
		Client:    client,
	}
}

// NewManager creates a new EC2InstanceManager using the default AWS config and client
func NewManager() (*EC2InstanceManager, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}

	client := ec2.NewFromConfig(cfg)

	mgr := NewManagerWithClient(client)

	return mgr, nil
}

// FetchAMIAge retrieves the AMI age of the given EC2 instance.
func (instance *EC2Instance) FetchAMIAge(client *ec2.Client) error {
	describeImagesInput := &ec2.DescribeImagesInput{
		ImageIds: []string{instance.AMI_ID},
	}

	imagesResp, err := client.DescribeImages(context.Background(), describeImagesInput)
	if err != nil {
		return fmt.Errorf("failed to describe image %s, %v", instance.AMI_ID, err)
	}

	if len(imagesResp.Images) == 0 {
		return fmt.Errorf("no image found with id %s", instance.AMI_ID)
	}

	amiCreationDate, err := time.Parse(time.RFC3339, *imagesResp.Images[0].CreationDate)
	if err != nil {
		return fmt.Errorf("failed to parse creation date of image %s, %v", instance.AMI_ID, err)
	}

	instance.AMI_Age = fmt.Sprintf("%d", int(time.Since(amiCreationDate).Hours()/24))

	return nil
}

// FetchAMIAge recursively calls FetchAMIAge on all Instances.
func (mgr *EC2InstanceManager) FetchAMIAge() error {
	for i := range mgr.Instances {
		err := mgr.Instances[i].FetchAMIAge(mgr.Client)
		if err != nil {
			return fmt.Errorf("unable to scan AMI Ages: %w", err)
		}
	}

	return nil
}

// FetchInstances connects to an AWS and fetches descriptions of all
// ec2 instances. These instances will be available in the Instances
// field.
func (mgr *EC2InstanceManager) FetchInstances(filter string) error {
	// empty in case of multiple runs
	mgr.Instances = []EC2Instance{}

	var filters []types.Filter
	if filter != "" {
		filter = "*" + filter + "*"

		filters = append(filters, types.Filter{
			Name:   aws.String("tag:Name"),
			Values: []string{filter},
		})
	}

	input := &ec2.DescribeInstancesInput{
		Filters: filters,
	}

	resp, err := mgr.Client.DescribeInstances(context.Background(), input)
	if err != nil {
		return fmt.Errorf("failed to describe instances: %w", err)
	}

	for _, reservation := range resp.Reservations {
		for _, instance := range reservation.Instances {
			var name string

			for _, tag := range instance.Tags {
				if *tag.Key == "Name" {
					name = *tag.Value
					break
				}
			}

			instanceStruct := EC2Instance{
				ID:          *instance.InstanceId,
				Name:        name,
				AMI_ID:      *instance.ImageId,
				InstanceAge: fmt.Sprintf("%d", int(time.Since(*instance.LaunchTime).Hours()/24)),
				AMI_Age:     "N/A", // not scanned initially
				IsSSM:       false, // This needs actual check
				Status:      string(instance.State.Name),
				PublicIP:    aws.ToString(instance.PublicIpAddress),
				PrivateIP:   aws.ToString(instance.PrivateIpAddress),
			}
			mgr.Instances = append(mgr.Instances, instanceStruct)
		}
	}

	return nil
}

//manager.Filter(ec2_instance.RunningInstances)

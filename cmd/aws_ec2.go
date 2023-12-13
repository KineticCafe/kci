package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type AWSInstance struct {
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

func GetInstances(filter string, includeAll bool) ([]AWSInstance, error) {

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return []AWSInstance{}, fmt.Errorf("unable to load SDK config: %w", err)
	}

	client := ec2.NewFromConfig(cfg)

	// Define filters based on flags
	var filters []types.Filter
	if filter != "" {
		filter = "*" + filter + "*"

		filters = append(filters, types.Filter{
			Name:   aws.String("tag:Name"),
			Values: []string{filter},
		})
	}

	if includeAll {
		filters = append(filters, types.Filter{
			Name:   aws.String("instance-state-name"),
			Values: []string{"pending", "running", "shutting-down", "terminated", "stopping"},
		})
	} else {
		filters = append(filters, types.Filter{
			Name:   aws.String("instance-state-name"),
			Values: []string{"pending", "running"},
		})
	}

	input := &ec2.DescribeInstancesInput{
		Filters: filters,
		// Populate filters
	}

	resp, err := client.DescribeInstances(context.Background(), input)
	if err != nil {
		return []AWSInstance{}, fmt.Errorf("failed to describe instances: %w", err)
	}

	var instances []AWSInstance
	for _, reservation := range resp.Reservations {
		for _, instance := range reservation.Instances {
			var name string
			for _, tag := range instance.Tags {
				if *tag.Key == "Name" {
					name = *tag.Value
					break
				}
			}

			// TODO add a filter for ages
			amiAge, err := getAMIAge(client, *instance.ImageId)
			if err != nil {
				return []AWSInstance{}, fmt.Errorf("unable to get AMI Age: %w", err)
			}

			instanceStruct := AWSInstance{
				ID:          *instance.InstanceId,
				Name:        name,
				AMI_ID:      *instance.ImageId,
				InstanceAge: fmt.Sprintf("%d", int(time.Since(*instance.LaunchTime).Hours()/24)),
				AMI_Age:     amiAge,
				IsSSM:       false, // This needs actual check
				Status:      string(instance.State.Name),
				PublicIP:    aws.ToString(instance.PublicIpAddress),
				PrivateIP:   aws.ToString(instance.PrivateIpAddress),
			}
			instances = append(instances, instanceStruct)
		}
	}

	return instances, nil
}

func getAMIAge(client *ec2.Client, amiID string) (string, error) {
	// Fetch AMI information
	describeImagesInput := &ec2.DescribeImagesInput{
		ImageIds: []string{amiID},
	}

	imagesResp, err := client.DescribeImages(context.Background(), describeImagesInput)
	if err != nil {
		return "", fmt.Errorf("failed to describe image %s, %v", amiID, err)
	}

	if len(imagesResp.Images) == 0 {
		return "", fmt.Errorf("no image found with id %s", amiID)
	}

	amiCreationDate, err := time.Parse(time.RFC3339, *imagesResp.Images[0].CreationDate)
	if err != nil {
		return "", fmt.Errorf("failed to parse creation date of image %s, %v", amiID, err)
	}

	amiAge := fmt.Sprintf("%d", int(time.Since(amiCreationDate).Hours()/24))

	return amiAge, nil
}

/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// TODO extract this when I have more than 1 command
type AWSInstance struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	AMI_ID      string `json:"ami_id"`
	InstanceAge string `json:"instance_age"`
	AMI_Age     string `json:"ami_age"`
	IsSSM       bool   `json:"is_ssm"`
	Status      string `json:"status"`
	PublicIP    string `json:"public_ip"`
	PrivateIP   string `jsoin:"private_ip"`
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

// instanceListCmd represents the instanceList command
var instanceListCmd = &cobra.Command{
	Use:   "list",
	Short: "list KCS instances",
	Long: `Lists KCS instances based on various criteria. Filters include:

- instance name
- SSM enabled
`,
	Run: func(cmd *cobra.Command, args []string) {

		filter, _ := cmd.Flags().GetString("filter")
		ssm, _ := cmd.Flags().GetBool("ssm")
		includeStopped, _ := cmd.Flags().GetBool("include-stopped")

		cfg, err := config.LoadDefaultConfig(context.Background())
		if err != nil {
			log.Fatalf("unable to load SDK config, %v", err)
		}

		client := ec2.NewFromConfig(cfg)

		// Define filters based on flags
		var filters []types.Filter
		if filter != "" {
			filters = append(filters, types.Filter{
				Name:   aws.String("tag:Name"),
				Values: []string{filter},
			})
		}
		//if !includeStopped {
		//	filters = append(filters, types.Filter{
		//		Name:   aws.String("instance-state-name"),
		//		Values: []string{"pending", "running", "shutting-down", "terminated", "stopping"},
		//	})
		//}
		input := &ec2.DescribeInstancesInput{
			Filters: filters,
			// Populate filters
		}

		resp, err := client.DescribeInstances(context.Background(), input)
		if err != nil {
			log.Fatalf("failed to describe instances, %v", err)
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
					log.Fatal(err)
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

		//spew.Dump(instances)
		// Now, instances slice contains details of all instances
		// Format this data as required (JSON, table, CSV)

		//sort.Slice(instances, func(i, j int) bool {
		//	return instances[i].Name < instances[j].Name
		//})

		sort.Slice(instances, func(i, j int) bool {
			iAge, _ := strconv.Atoi(instances[i].InstanceAge)
			jAge, _ := strconv.Atoi(instances[j].InstanceAge)
			return iAge < jAge
		})

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Name", "AMI ID", "Instance Age", "AMI Age", "Is SSM", "Status", "PublicIP", "PrivateIP"})

		for _, instance := range instances {
			table.Append([]string{
				instance.ID,
				instance.Name,
				instance.AMI_ID,
				instance.InstanceAge,
				instance.AMI_Age,
				strconv.FormatBool(instance.IsSSM),
				instance.Status,
				instance.PublicIP,
				instance.PrivateIP,
			})
		}

		table.Render()
	},
}

func init() {
	instanceCmd.AddCommand(instanceListCmd)
	instanceListCmd.Flags().String("filter", "", "Filter instances by name")
	instanceListCmd.Flags().Bool("ssm", false, "Return only instances in SSM")
	instanceListCmd.Flags().Bool("include-stopped", false, "Include stopped instances in the list")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// instanceListCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// instanceListCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

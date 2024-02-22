package cmd

import (
	"log"
	"os"
	"sort"
	"strconv"

	"github.com/KineticCommerce/kci/ec2_instance"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var instanceAgingCmd = &cobra.Command{
	Use:   "aging",
	Short: "list KCS instances that are a little long in the tooth",
	Run: func(cmd *cobra.Command, args []string) {
		filter, _ := cmd.Flags().GetString("filter")
		includeAll, _ := cmd.Flags().GetBool("all")

		// Get with main filters
		manager, err := ec2_instance.NewManager()
		if err != nil {
			log.Fatal(err)
		}
		err = manager.FetchInstances(filter)
		if err != nil {
			log.Fatal(err)
		}

		if !includeAll {
			manager.Filter(ec2_instance.IsRunningFilter)
		}

		// Filter
		err = manager.FetchAMIAge()
		if err != nil {
			log.Fatal(err)
		}

		manager.Filter(func(instance ec2_instance.EC2Instance) bool {
			iAge, _ := strconv.Atoi(instance.InstanceAge)
			aAge, _ := strconv.Atoi(instance.AMI_Age)

			return iAge > 90 || aAge > 90
		})

		// Display
		sort.Slice(manager.Instances, func(i, j int) bool {
			iAge, _ := strconv.Atoi(manager.Instances[i].InstanceAge)
			jAge, _ := strconv.Atoi(manager.Instances[j].InstanceAge)
			return iAge < jAge
		})

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "ID", "AMI ID", "Instance Age", "AMI Age", "Status"})

		for _, instance := range manager.Instances {
			table.Append([]string{
				instance.Name,
				instance.ID,
				instance.AMI_ID,
				instance.InstanceAge,
				instance.AMI_Age,
				instance.Status,
			})
		}

		table.Render()
	},
}

func init() {
	instanceCmd.AddCommand(instanceAgingCmd)
	instanceAgingCmd.Flags().StringP("filter", "f", "", "Filter instances by name")
	instanceAgingCmd.Flags().BoolP("all", "a", false, "Include all statuses in the list")
}

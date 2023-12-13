package cmd

import (
	"log"
	"os"
	"sort"

	"github.com/KineticCommerce/kci/ec2_instance"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// instanceListCmd represents the instanceList command
var instanceListCmd = &cobra.Command{
	Use:   "list",
	Short: "list KCS instances",
	Run: func(cmd *cobra.Command, args []string) {
		filter, _ := cmd.Flags().GetString("filter")
		includeAll, _ := cmd.Flags().GetBool("all")

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

		sort.Slice(manager.Instances, func(i, j int) bool {
			return manager.Instances[i].Name < manager.Instances[j].Name
		})

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "ID", "Instance Age", "Status", "PublicIP", "PrivateIP"})

		for _, instance := range manager.Instances {
			table.Append([]string{
				instance.Name,
				instance.ID,
				//instance.AMI_ID,
				instance.InstanceAge,
				//instance.AMI_Age,
				//strconv.FormatBool(instance.IsSSM),
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
	instanceListCmd.Flags().StringP("filter", "f", "", "Filter instances by name")
	instanceListCmd.Flags().BoolP("all", "a", false, "Include all statuses in the list")
	instanceListCmd.Flags().Bool("ssm", false, "Only show instances with SSM enabled")
	instanceListCmd.Flags().Bool("no-ssm", false, "Only show instances without SSM enabled")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// instanceListCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// instanceListCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

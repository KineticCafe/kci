/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"
	"strconv"

	"github.com/KineticCommerce/kci/ec2_instance"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	_ "net/http/pprof"
)

// scanCmd represents the scan command
var instanceScanCmd = &cobra.Command{
	Use:   "scan",
	Short: "scan instances for OS and reboot status",
	Run: func(cmd *cobra.Command, args []string) {
		filter, _ := cmd.Flags().GetString("filter")
		includeAll, _ := cmd.Flags().GetBool("all")

		// local filters and flags
		jump, _ := cmd.Flags().GetString("jump")
		jumpuser, _ := cmd.Flags().GetString("jumpuser")
		rebootOnly, _ := cmd.Flags().GetBool("reboot-only")

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

		// Filter and Calculate
		// We pre-filter here because this is expensive and
		// we are sure that we do not want to scan these on this pass.
		manager.Filter(func(instance ec2_instance.EC2Instance) bool {
			return len(instance.PublicIP) == 0
		})

		// TODO move to package
		//err = ScanInstances(jump, jumpuser, manager.Instances)
		err = manager.JumpScan(jump, jumpuser)
		if err != nil {
			log.Fatal(err)
		}

		if rebootOnly {
			manager.Filter(func(instance ec2_instance.EC2Instance) bool {
				return instance.RebootRequired == "reboot required"
			})
		}
		// Display
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "ID", "Instance Age", "Uptime", "Reboot", "Updates", "OS", "Private IP"})

		for _, instance := range manager.Instances {
			updates := strconv.Itoa(instance.SecurityUpdates)

			table.Append([]string{
				instance.Name,
				instance.ID,
				instance.InstanceAge,
				instance.Uptime,
				instance.RebootRequired,
				updates,
				instance.OsVersion,
				instance.PrivateIP,
			})
		}

		table.Render()
	},
}

func init() {
	instanceCmd.AddCommand(instanceScanCmd)
	instanceScanCmd.Flags().StringP("filter", "f", "", "Filter instances by name")
	instanceScanCmd.Flags().StringP("jump", "j", "", "jumpbox server address")
	instanceScanCmd.Flags().StringP("jumpuser", "u", "", "jumpbox user")
	instanceScanCmd.Flags().BoolP("all", "a", false, "Include all statuses in the list")
	instanceScanCmd.Flags().Bool("reboot-only", false, "Display only servers that need a reboot")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

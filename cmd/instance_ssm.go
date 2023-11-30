/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
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

// instanceAgingCmd represents the instanceAging command
var instanceSSMCmd = &cobra.Command{
	Use:   "ssm",
	Short: "list KCS instances that are managed by SSM",
	Run:   listSSMCommand,
}

func listSSMCommand(cmd *cobra.Command, args []string) {
	filter, _ := cmd.Flags().GetString("filter")
	disabled, _ := cmd.Flags().GetBool("disabled")

	// Get with main filters
	manager, err := ec2_instance.NewManager()
	if err != nil {
		log.Fatal(err)
	}
	err = manager.FetchInstances(filter)
	if err != nil {
		log.Fatal(err)
	}

	manager.Filter(ec2_instance.IsRunningFilter)

	err = manager.FetchSSMDetails()
	if err != nil {
		log.Fatal(err)
	}

	manager.Filter(func(i ec2_instance.EC2Instance) bool {
		if disabled {
			return !i.IsSSM
		} else {
			return i.IsSSM
		}
	})

	// Display
	sort.Slice(manager.Instances, func(i, j int) bool {
		iAge, _ := strconv.Atoi(manager.Instances[i].InstanceAge)
		jAge, _ := strconv.Atoi(manager.Instances[j].InstanceAge)
		return iAge < jAge
	})

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "ID", "SSM Enabled", "Status"})

	for _, instance := range manager.Instances {
		table.Append([]string{
			instance.Name,
			instance.ID,
			strconv.FormatBool(instance.IsSSM),
			instance.Status,
		})
	}

	table.Render()
}

func init() {
	instanceCmd.AddCommand(instanceSSMCmd)
	instanceSSMCmd.Flags().StringP("filter", "f", "", "Filter instances by name")
	instanceSSMCmd.Flags().Bool("disabled", false, "Display SSM disabled instead")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// instanceAgingCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// instanceAgingCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

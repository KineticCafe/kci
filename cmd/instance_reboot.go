/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/KineticCommerce/kci/ec2_instance"
	"github.com/spf13/cobra"
)

// instanceRebootCmd represents the instanceReboot command
var instanceRebootCmd = &cobra.Command{
	Use:   "reboot",
	Short: "reboot an instance",
	Run: func(cmd *cobra.Command, args []string) {

		instanceID, _ := cmd.Flags().GetString("instance-id")

		if instanceID == "" {
			log.Fatalf("instance-id is required")
		}

		err := ec2_instance.Reboot(instanceID)
		if err != nil {
			log.Fatalf("Unable to reboot instance %q, %v", instanceID, err)
		}

		log.Printf("Successfully requested reboot for instance %q", instanceID)
	},
}

func init() {
	instanceCmd.AddCommand(instanceRebootCmd)

	instanceRebootCmd.Flags().StringP("instance-id", "i", "", "the instance id to reboot")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// instanceRebootCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// instanceRebootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

package cmd

import (
	"log"

	"github.com/KineticCommerce/kci/ec2_instance"
	"github.com/spf13/cobra"
)

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
}

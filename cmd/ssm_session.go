package cmd

import (
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var ssmSessionCmd = &cobra.Command{
	Use:   "session",
	Short: "Start an SSM session for a given instance. AWS CLI required.",
	Run: func(cmd *cobra.Command, args []string) {
		instanceID, _ := cmd.Flags().GetString("instance-id")
		if instanceID == "" {
			log.Fatalf("instance-id is a required parameter")
		}

		c := exec.Command("aws", "ssm", "start-session", "--target", instanceID)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Stdin = os.Stdin

		// Run the command. This will block until the session is finished.
		err := c.Run()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	// this SHOULD work as a simple alias
	ssmCmd.AddCommand(ssmSessionCmd)
	ssmSessionCmd.Flags().StringP("instance-id", "i", "", "The instance to connect to - required")
}

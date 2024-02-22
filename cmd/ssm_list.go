package cmd

import "github.com/spf13/cobra"

var ssmListCmd = &cobra.Command{
	Use:   "list",
	Short: "list KCS instances with SSM enabled",
	Run:   listSSMCommand,
}

func init() {
	// this SHOULD work as a simple alias
	ssmCmd.AddCommand(ssmListCmd)
	ssmListCmd.Flags().StringP("filter", "f", "", "Filter instances by name")
	ssmListCmd.Flags().Bool("disabled", false, "Display SSM disabled instead")
}

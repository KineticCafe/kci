package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var sysinfoCmd = &cobra.Command{
	Use:   "sysinfo",
	Short: "Subcommands for querying KCS system information",
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(sysinfoCmd)
}

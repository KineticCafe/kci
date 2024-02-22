package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var ssmCmd = &cobra.Command{
	Use:   "ssm",
	Short: "Subcommands for interacting with KCS instances over SSM",
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(ssmCmd)
}

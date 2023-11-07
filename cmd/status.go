/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// instanceCmd represents the instance command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Subcommands for querying KCS service status",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("status called")
		fmt.Println(environment)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// instanceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// instanceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

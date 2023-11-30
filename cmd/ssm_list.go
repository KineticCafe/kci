/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// instanceListCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// instanceListCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

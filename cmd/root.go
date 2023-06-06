/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	debug             bool
	verbose           bool
	environment       string
	validEnvironments = []string{"dit", "stage", "prod"}
)

func isValidEnvironment(env string) bool {
	for _, b := range validEnvironments {
		if b == env {
			return true
		}
	}
	return false
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kcs-infra",
	Short: "KCS infrastructure and application management",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if !isValidEnvironment(environment) {
			return fmt.Errorf("Invalid environment: '%v'. Environment can be 'dit', 'stage', or 'prod'", environment)
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kcs-infra.yaml)")
	rootCmd.PersistentFlags().StringVarP(&environment, "environment", "e", "dit", "Set the environment. Can be 'dit', 'stage', or 'prod'")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "debug output")
	//rootCmd.PersistentFlags().SetAnnotation("environment", cobra.BashCompOneOf, validEnvironements)

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

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
	validEnvironments = []string{"dit", "stage", "prod", "prod-eu"}
)

func isValidEnvironment(env string) bool {
	for _, b := range validEnvironments {
		if b == env {
			return true
		}
	}
	return false
}

var rootCmd = &cobra.Command{
	Use:   "kci",
	Short: "KCS infrastructure management and reporting",
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

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&environment, "environment", "e", "dit", "Set the environment. Can be 'dit', 'stage', or 'prod'")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "debug output")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

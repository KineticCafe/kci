package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// BuildTime should be set using ldflags during the build process.
// e.g. -ldflags "-X 'github.com/KineticCommerce/kci/cmd.BuildTime=$(date)'"
var BuildTime = "not set"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version and build information",
	Run: func(cmd *cobra.Command, args []string) {
		// Just going to use the build time for the version for the time
		// being. That's all I really need right now.
		fmt.Println(rootCmd.Use, "-", rootCmd.Short)
		fmt.Println("build: ", BuildTime)
		fmt.Println("repo:   https://github.com/KineticCommerce/kci")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

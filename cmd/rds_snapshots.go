package cmd

import (
	"log"
	"os"

	"github.com/KineticCommerce/kci/database"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var rdsSnapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "list snapshots for a database",
	Run: func(cmd *cobra.Command, args []string) {
		identifier, _ := cmd.Flags().GetString("identifier")

		manager, err := database.NewManager()
		if err != nil {
			log.Fatalf("unable to load database manager %v", err)
		}

		snapshots, err := manager.FetchSnapshots(identifier)
		if err != nil {
			log.Fatalf("unlable to load database %s: %v", identifier, err)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Created At"})

		for _, snapshot := range snapshots {

			table.Append([]string{
				snapshot.ID,
				snapshot.Created.Format("2006-01-02 15:04:05"),
			})
		}

		table.Render()
	},
}

func init() {
	rdsCmd.AddCommand(rdsSnapshotCmd)
	rdsSnapshotCmd.Flags().StringP("identifier", "i", "", "database identifier")
}

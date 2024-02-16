package cmd

import (
	"log"
	"os"
	"strconv"

	"github.com/KineticCommerce/kci/database"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var rdsListCmd = &cobra.Command{
	Use:   "list",
	Short: "list KCS RDS databases",
	Run: func(cmd *cobra.Command, args []string) {
		filter, _ := cmd.Flags().GetString("filter")
		manager, err := database.NewManager()
		if err != nil {
			log.Fatalf("unable to load SDK config, %v", err)
		}

		err = manager.Fetch(filter)
		if err != nil {
			log.Fatalf("unlable to load databases: %v", err)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Multi AZ", "Latest Snapshot ID", "Snapshot Count"})

		for _, db := range manager.Databases {

			table.Append([]string{
				db.ID,
				strconv.FormatBool(db.MultiAZ),
				db.LatestSnapshotID(),
				strconv.Itoa(len(db.Snapshots)),
			})
		}

		table.Render()
	},
}

func init() {
	rdsCmd.AddCommand(rdsListCmd)

	rdsListCmd.Flags().StringP("filter", "f", "", "Filter databases by name")
}

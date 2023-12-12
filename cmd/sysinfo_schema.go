/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/

package cmd

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

type Schema struct {
	Change    string `json:"change"`
	ChangeID  string `json:"change_id"`
	PlannedAt string `json:"planned_at"`
	Hashref   string `json:"script_hash"`
}

// TODO possibly remove the hardcode - not necessary right now
type SchemaResponse struct {
	KiehlsCAS Schema `json:"kinetic-cas-kiehls-schema"`
	Core      Schema `json:"kinetic-platform-schema"`
}

var sysinfoSchemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "display current sqitch migration for deployed services",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO Move to function when used more than once...

		urls := map[string]string{
			"dit":     "https://kcs-dev.kineticcommercetech.io/status/sqitch",
			"staging": "https://kcs-staging.kineticcommercetech.io/status/sqitch",
			"prod":    "https://kcs.kineticcommerce.io/status/sqitch",
			"prod-eu": "https://kcs-prod-eu-platform.kineticcommerce.io/status/sqitch",
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Env", "Core Change", "Core Planned At", "Kiehls Change", "Kiehls Planned At"})

		// TODO probably a good idea to move this out of the display loop
		// TODO with multiple DBs this is not as clean as it could be
		for key, url := range urls {
			resp, err := http.Get(url)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				log.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}

			var result SchemaResponse
			err = json.Unmarshal(body, &result)
			if err != nil {
				log.Fatal(err)
			}

			table.Append([]string{
				key,
				result.Core.ChangeID,
				result.Core.PlannedAt,
				result.KiehlsCAS.ChangeID,
				result.KiehlsCAS.PlannedAt,
			})
		}

		table.Render()
	},
}

func init() {
	sysinfoCmd.AddCommand(sysinfoSchemaCmd)
	//instanceSSMCmd.Flags().string("filter", "", "Filter instances by name")
	//instanceSSMCmd.Flags().Bool("disabled", false, "Display SSM disabled instead")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// instanceAgingCmd.PersistentFlags().string("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// instanceAgingCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

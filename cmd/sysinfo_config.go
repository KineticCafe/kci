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

type Config struct {
	Hashref   string `json:"hashref"`
	Timestamp string `json:"timestamp"`
}

var sysinfoConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "display current config package",
	Run: func(cmd *cobra.Command, args []string) {
		/*
			TODO Move to function when used more than once...
		*/

		urls := map[string]string{
			"dit":     "https://kcs-dev.kineticcommercetech.io/status/config",
			"staging": "https://kcs-staging.kineticcommercetech.io/status/config",
			"prod":    "https://kcs.kineticcommerce.io/status/config",
			"prod-eu": "https://kcs-prod-eu-platform.kineticcommerce.io/status/config",
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Env", "Hashref", "Timestamp"})

		// TODO probably a good idea to move this out of the display loop
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

			var result Config
			err = json.Unmarshal(body, &result)

			if err != nil {
				log.Fatal(err)
			}

			table.Append([]string{
				key,
				result.Hashref,
				result.Timestamp,
			})
		}

		table.Render()
	},
}

func init() {
	sysinfoCmd.AddCommand(sysinfoConfigCmd)
	//instanceSSMCmd.Flags().String("filter", "", "Filter instances by name")
	//instanceSSMCmd.Flags().Bool("disabled", false, "Display SSM disabled instead")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// instanceAgingCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// instanceAgingCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

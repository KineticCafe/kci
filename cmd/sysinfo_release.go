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

// Release represents a platform release as returned by the `status/release` endpoint.
type Release struct {
	Elixir    interface{} `json:"elixir"`
	Hashref   string      `json:"hashref"`
	Name      string      `json:"name"`
	Repo      interface{} `json:"repo"`
	Timestamp string      `json:"timestamp"`
}

var sysinfoReleaseCmd = &cobra.Command{
	Use:   "release",
	Short: "display current release information",
	Run: func(cmd *cobra.Command, args []string) {
		/*
			TODO Move to function when used more than once...
		*/

		urls := map[string]string{
			"dit":     "https://kcs-dev.kineticcommercetech.io/status/release",
			"staging": "https://kcs-staging.kineticcommercetech.io/status/release",
			"prod":    "https://kcs.kineticcommerce.io/status/release",
			"prod-eu": "https://kcs-prod-eu-platform.kineticcommerce.io/status/release",
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

			var result map[string]Release
			err = json.Unmarshal(body, &result)
			if err != nil {
				log.Fatal(err)
			}

			table.Append([]string{
				key,
				result["package"].Hashref,
				result["package"].Timestamp,
			})
		}

		table.Render()
	},
}

func init() {
	sysinfoCmd.AddCommand(sysinfoReleaseCmd)
}

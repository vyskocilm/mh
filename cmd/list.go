package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"
)

var (
	printJSON bool
)

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVarP(&printJSON, "json", "", false, "print output as JSON")
}

var listCmd = &cobra.Command {
	Use:   "list",
	Short: "list entries in hosts",
	Long:  `list prints all managed entries from hosts file`,
	Args: func(cmd *cobra.Command, args []string) error {
        err := parseArgs(cmd, args)
        if err != nil {
            return err
        }
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		c := cfg.HTTPClient()
		req, err := http.NewRequest(
			"GET",
			"http://localhost:1234/v1/e",
			nil,
		)
		if err != nil {
			exitOnErr(err)
		}
		resp, err := c.Do(req)
		if err != nil {
			exitOnErr(err)
		}
		if resp.StatusCode != http.StatusOK {
			exitOnResp(resp)
		}
		// read reply
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			exitOnErr(err)
		}
		err = resp.Body.Close()
		if err != nil {
			exitOnErr(err)
		}
		if !printJSON {
			// unmarshal
			var entries map[string][]e
			err = json.Unmarshal(b, &entries)
			if err != nil {
				exitOnErr(err)
			}
            for group, es := range entries {
                fmt.Printf("%s: \n", group)
                for _, foo := range es {
                    fmt.Printf("    %s\t%s\n", foo.IP, foo.Name)
                }
            }
		} else {
			fmt.Printf("%s", string(b))
		}
	},
}

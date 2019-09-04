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

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list entries in hosts",
	Long:  `list prints all managed entries from hosts file`,
	Run: func(cmd *cobra.Command, args []string) {

		c := http.Client{}
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
			var list []e
			err = json.Unmarshal(b, &list)
			if err != nil {
				exitOnErr(err)
			}
			for _, foo := range list {
				fmt.Printf("%s\t%s\n", foo.IP, foo.Name)
			}
		} else {
			fmt.Printf("%s", string(b))
		}
	},
}

package cmd

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(delCmd)
    addGroupFlag(delCmd)
}

var delCmd = &cobra.Command{
	Use:   "del",
	Short: "delete entries from hosts",
	Long:  `delete first matching IP or name from hosts, supports more arguments`,
	Args: func(cmd *cobra.Command, args []string) error {
        err := parseArgs(cmd, args)
        applyGroupEnv()
        if err != nil {
            return err
        }
		if len(args) < 1 {
			return errors.New("Command del requires at least one `ip' or `name' argument")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		c := cfg.HTTPClient()
		for i := range args {
			req, err := http.NewRequest(
				"DELETE",
				fmt.Sprintf("http://localhost:1234/v1/e/%s/%s", groupVar, args[i]),
				nil,
			)
			if err != nil {
				exitOnErr(err)
			}
			resp, err := c.Do(req)
			if err != nil {
				exitOnErr(err)
			}
			if resp.StatusCode != http.StatusNoContent {
				exitOnResp(resp)
			}
		}
	},
}

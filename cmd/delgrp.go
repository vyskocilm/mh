package cmd

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(delGrpCmd)
}

var delGrpCmd = &cobra.Command{
	Use:   "delgrp",
	Short: "delete all entries from the group from hosts",
	Long:  `delete all groups from the hosts, supports more arguments or MH_GROUP variable`,
	Args: func(cmd *cobra.Command, args []string) error {
        err := parseArgs(cmd, args)
        applyGroupEnv()
        if err != nil {
            return err
        }
		if len(args) < 1 || groupVar != defaultGROUP {
			return errors.New("Command delgrp requires at least one `group'argument, or MH_GROUP environment variable")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

        if groupVar != defaultGROUP {
            args = append(args, defaultGROUP)
        }

		c := cfg.HTTPClient()
		for i := range args {
			req, err := http.NewRequest(
				"DELETE",
				fmt.Sprintf("http://localhost:1234/v1/g/%s", args[i]),
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

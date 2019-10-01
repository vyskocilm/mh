package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(addCmd)
	addGroupFlag(addCmd)
}

// print given error and call os.Exit
func exitOnErr(err error) {
	fmt.Printf("%s\n", err.Error())
	os.Exit(1)
}

// print returned reply and call os.Exit
func exitOnResp(resp *http.Response) {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		exitOnErr(err)
	}
	err = resp.Body.Close()
	if err != nil {
		exitOnErr(err)
	}
	exitOnErr(
		fmt.Errorf("Received %s: %s", resp.Status, string(b)))
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add entry to hosts",
	Long:  `add new ip name entry to hosts file`,
	Args: func(cmd *cobra.Command, args []string) error {
		err := parseArgs(cmd, args)
		applyGroupEnv()
		if err != nil {
			return err
		}
		if len(args) < 2 {
			return errors.New("Command add requires `ip' `name' arguments")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		c := cfg.HTTPClient()
		req, err := http.NewRequest(
			"PUT",
			fmt.Sprintf("http://localhost:1234/v1/e/%s/%s/%s", groupVar, args[0], args[1]),
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

	},
}

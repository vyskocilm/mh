package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mh",
	Short: "mh let you dynamically manage /etc/hosts",
	Long: `Easy to use way of managing /etc/hosts
        with a simple and clean command line interface and garbage collection`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

// Execute run mh command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

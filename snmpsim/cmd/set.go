package cmd

import (
	"fmt"
	setsubcommands "snmpsim-cli-manager/snmpsim/cmd/setSubcommands"

	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Sets a certain property",
	Long:  `Sets a certain property of the cli-management tool`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("set called")
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
	setCmd.AddCommand(setsubcommands.SetConfigCmd)
	setCmd.AddCommand(setsubcommands.SetEnvConfigCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

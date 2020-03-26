package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"snmpsim-cli-manager/snmpsim/cmd/addSubcommands"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add <sub-componentToComponent>",
	Args:  cobra.ExactArgs(0),
	Short: "Adds a sub-component to a component",
	Long: `Adds a certain sub-component to a component.

A sub-component can only be added to its related main-component:
	- An user and an endpoint can only be added to an engine
	- An engine can only be added to an agent
	- An agent can only be added to a lab`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd.Help())
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.AddCommand(addsubcommands.AgentToLabCmd)
	addCmd.AddCommand(addsubcommands.EndpointToEngineCmd)
	addCmd.AddCommand(addsubcommands.EngineToAgentCmd)
	addCmd.AddCommand(addsubcommands.UserToEngineCmd)
}

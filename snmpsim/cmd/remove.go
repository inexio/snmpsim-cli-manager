package cmd

import (
	"fmt"
	"snmpsim-cli-manager/snmpsim/cmd/removeSubcommands"

	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove <sub-componentFromComponent>",
	Args:  cobra.ExactArgs(0),
	Short: "Removes a sub-component from a component",
	Long: `Removes a certain sub-component from a component.

A sub-component can only be removed from its related main-component:
	- An agent can only be removed from a lab
	- An engine can only be removed from an agent
	- An user and an endpoint can only be removed from an engine`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd.Help())
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
	removeCmd.AddCommand(removesubcommands.AgentFromLabCmd)
	removeCmd.AddCommand(removesubcommands.EndpointFromEngineCmd)
	removeCmd.AddCommand(removesubcommands.EngineFromAgentCmd)
	removeCmd.AddCommand(removesubcommands.UserFromEngineCmd)
}

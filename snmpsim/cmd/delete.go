package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"snmpsim-cli-manager/snmpsim/cmd/deleteSubcommands"
)

// deleteCmd represents the deleteSubcommands command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Args:  cobra.ExactArgs(0),
	Short: "Deletes the component with the given id",
	Long:  `Completely deletes a component via the given id.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd.Help())
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.AddCommand(deletesubcommands.DeleteAgentCmd)
	deleteCmd.AddCommand(deletesubcommands.DeleteEndpointCmd)
	deleteCmd.AddCommand(deletesubcommands.DeleteEngineCmd)
	deleteCmd.AddCommand(deletesubcommands.DeleteLabCmd)
	deleteCmd.AddCommand(deletesubcommands.DeleteRecordFileCmd)
	deleteCmd.AddCommand(deletesubcommands.DeleteTagCmd)
	deleteCmd.AddCommand(deletesubcommands.DeleteUserCmd)
	deleteCmd.AddCommand(deletesubcommands.DeleteTaggedObjectsCmd)

}

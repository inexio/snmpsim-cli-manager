package cmd

import (
	"fmt"
	"snmpsim-cli-manager/snmpsim/cmd/getSubcommands"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Args:  cobra.ExactArgs(0),
	Short: "Returns data about a component",
	Long:  `Returns detailed information about the given component.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd.Help())
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.PersistentFlags().String("format", "human-readable", "Use this flag to set the output format")
	getCmd.PersistentFlags().Int("depth", 3, "Use this flag to set the depth of human-readable output")
	getCmd.PersistentFlags().BoolP("pretty", "p", false, "Use this flag to prettify the output if xml or json is set as output format")
	getCmd.AddCommand(getsubcommands.GetAgentCmd)
	getCmd.AddCommand(getsubcommands.GetAgentsCmd)
	getCmd.AddCommand(getsubcommands.GetEndpointCmd)
	getCmd.AddCommand(getsubcommands.GetEndpointsCmd)
	getCmd.AddCommand(getsubcommands.GetEngineCmd)
	getCmd.AddCommand(getsubcommands.GetEnginesCmd)
	getCmd.AddCommand(getsubcommands.GetFilterValuesCmd)
	getCmd.AddCommand(getsubcommands.GetLabCmd)
	getCmd.AddCommand(getsubcommands.GetLabsCmd)
	getCmd.AddCommand(getsubcommands.MessageFiltersCmd)
	getCmd.AddCommand(getsubcommands.MessagesCmd)
	getCmd.AddCommand(getsubcommands.PacketFiltersCmd)
	getCmd.AddCommand(getsubcommands.PacketsCmd)
	getCmd.AddCommand(getsubcommands.ProcessCmd)
	getCmd.AddCommand(getsubcommands.ProcessesCmd)
	getCmd.AddCommand(getsubcommands.GetRecordFilesCmd)
	getCmd.AddCommand(getsubcommands.GetRecordFileCmd)
	getCmd.AddCommand(getsubcommands.GetTagCmd)
	getCmd.AddCommand(getsubcommands.GetTagsCmd)
	getCmd.AddCommand(getsubcommands.GetUserCmd)
	getCmd.AddCommand(getsubcommands.GetUsersCmd)
}

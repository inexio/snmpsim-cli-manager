package cmd

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
	"snmpsim-cli-manager/snmpsim/cmd/untagSubcommands"
)

// untagCmd represents the untag command
var untagCmd = &cobra.Command{
	Use:   "untag <component> <component-id> --tag <tag-id>",
	Args:  cobra.ExactArgs(0),
	Short: "Untags a given component",
	Long:  `Removes the tag with the given tag-id from the component with the given component-id.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd.Help())
	},
}

func init() {
	rootCmd.AddCommand(untagCmd)
	untagCmd.AddCommand(untagsubcommands.UntagLabCmd)
	untagCmd.AddCommand(untagsubcommands.UntagEngineCmd)
	untagCmd.AddCommand(untagsubcommands.UntagAgentCmd)
	untagCmd.AddCommand(untagsubcommands.UntagEndpointCmd)
	untagCmd.AddCommand(untagsubcommands.UntagUserCmd)
	untagCmd.PersistentFlags().Int("tag", 0, "The id of the tag")
	err := untagCmd.MarkPersistentFlagRequired("tag")
	if err != nil {
		log.Error().
			Err(err).
			Msg("Could not mark 'tag' flag required")
		os.Exit(1)
	}
}

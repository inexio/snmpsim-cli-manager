package cmd

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"snmpsim-cli-manager/snmpsim/cmd/tagSubcommands"

	"github.com/spf13/cobra"
)

// tagCmd represents the tag command
var tagCmd = &cobra.Command{
	Use:   "tag <component> <component-id> --tag <tag-id>",
	Args:  cobra.ExactArgs(0),
	Short: "Tags a component",
	Long:  `Tags the component with the given component-id with the given tag-id`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd.Help())
	},
}

func init() {
	rootCmd.AddCommand(tagCmd)
	tagCmd.AddCommand(tagsubcommands.TagLabCmd)
	tagCmd.AddCommand(tagsubcommands.TagEngineCmd)
	tagCmd.AddCommand(tagsubcommands.TagAgentCmd)
	tagCmd.AddCommand(tagsubcommands.TagEndpointCmd)
	tagCmd.AddCommand(tagsubcommands.TagUserCmd)
	tagCmd.PersistentFlags().Int("tag", 0, "The id of the tag")
	err := tagCmd.MarkPersistentFlagRequired("tag")
	if err != nil {
		log.Error().
			Msg("Could not mark 'tag' flag required")
		os.Exit(1)
	}
}

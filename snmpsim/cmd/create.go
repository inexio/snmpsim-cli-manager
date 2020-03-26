package cmd

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
	"snmpsim-cli-manager/snmpsim/cmd/createSubcommands"
)

// createCmd represents the createSubcommands command
var createCmd = &cobra.Command{
	Use:   "create",
	Args:  cobra.ExactArgs(0),
	Short: "Creates a component",
	Long: `Creates a component in the snmpsim environment.
Components can be created with or without a tag.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd.Help())
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.PersistentFlags().String("name", "", "Name of the component")
	err := createCmd.MarkPersistentFlagRequired("name")
	if err != nil {
		log.Error().
			Msg("Could not mark 'name' flag required")
		os.Exit(1)
	}
	createCmd.AddCommand(createsubcommands.CreateAgentCmd)
	createCmd.AddCommand(createsubcommands.CreateEndpointCmd)
	createCmd.AddCommand(createsubcommands.CreateEngineCmd)
	createCmd.AddCommand(createsubcommands.CreateLabCmd)
	createCmd.AddCommand(createsubcommands.CreateTagCmd)
	createCmd.AddCommand(createsubcommands.CreateUserCmd)
}

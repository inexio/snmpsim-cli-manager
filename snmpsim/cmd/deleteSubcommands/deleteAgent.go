package deletesubcommands

import (
	"fmt"
	snmpsimclient "github.com/inexio/snmpsim-restapi-go-client"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strconv"
)

// DeleteAgentCmd represents the deleteAgent command
var DeleteAgentCmd = &cobra.Command{
	Use:   "agent <id>",
	Args:  cobra.ExactArgs(1),
	Short: "Deletes an agent",
	Long:  `Deletes the agent with the given agent-id`,
	Run: func(deleteAgentCmd *cobra.Command, args []string) {
		//Load the client data from the config
		baseURL := viper.GetString("mgmt.http.baseURL")
		username := viper.GetString("mgmt.http.authUsername")
		password := viper.GetString("mgmt.http.authPassword")

		//Create a new client
		client, err := snmpsimclient.NewManagementClient(baseURL)
		if err != nil {
			log.Error().
				Msg("Error while creating management client")
			os.Exit(1)
		}
		err = client.SetUsernameAndPassword(username, password)
		if err != nil {
			log.Error().
				Msg("Error while setting username and password")
			os.Exit(1)
		}

		//Read in the agent-id
		agentID, err := strconv.Atoi(args[0])
		if err != nil {
			log.Error().
				Msg("Error during conversion of agent-id from string to integer")
			os.Exit(1)
		}

		//Delete the agent
		err = client.DeleteAgent(agentID)
		if err != nil {
			log.Error().
				Msg("Error while deleting agent")
			os.Exit(1)
		}

		fmt.Println("Agent", args[0], "has been deleted successfully.")
	},
}

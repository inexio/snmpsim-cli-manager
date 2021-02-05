package removesubcommands

import (
	"fmt"
	"os"

	snmpsimclient "github.com/inexio/snmpsim-restapi-go-client"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// EngineFromAgentCmd represents the engineFromAgent command
var EngineFromAgentCmd = &cobra.Command{
	Use:   "engine-from-agent",
	Args:  cobra.ExactArgs(0),
	Short: "Removes an engine from an agent",
	Long:  `Removes the engine with the given engine-id from the agent with the give agent-id`,
	Run: func(cmd *cobra.Command, args []string) {
		//Load the client data from the config
		baseURL := viper.GetString("mgmt.http.baseURL")
		username := viper.GetString("mgmt.http.authUsername")
		password := viper.GetString("mgmt.http.authPassword")

		//Create a new client
		client, err := snmpsimclient.NewManagementClient(baseURL)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error while creating management client")
			os.Exit(1)
		}
		if username != "" && password != "" {
			err = client.SetUsernameAndPassword(username, password)
			if err != nil {
				log.Error().
					Err(err).
					Msg("Error while setting username and password")
				os.Exit(1)
			}
		}

		//Read in the engine-id
		engineID, err := cmd.Flags().GetInt("engine")
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error while retrieving engineID")
			os.Exit(1)
		}

		//Read in the agent-id
		agentID, err := cmd.Flags().GetInt("agent")
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error while retrieving agentID")
			os.Exit(1)
		}

		//Remove the engine from the agent
		err = client.RemoveEngineFromAgent(engineID, agentID)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error while removing the engine from the agent")
			os.Exit(1)
		}
		fmt.Println("Engine", engineID, "has been removed from agent", agentID)
	},
}

func init() {
	//Set engine flag
	EngineFromAgentCmd.Flags().Int("engine", 0, "Id of the engine the will be remove from the agent")
	err := EngineFromAgentCmd.MarkFlagRequired("engine")
	if err != nil {
		log.Error().
			Err(err).
			Msg("Could not mark 'engine' flag required")
		os.Exit(1)
	}

	//Set agent flag
	EngineFromAgentCmd.Flags().Int("agent", 0, "Id of the agent the engine will be removed from")
	err = EngineFromAgentCmd.MarkFlagRequired("agent")
	if err != nil {
		log.Error().
			Err(err).
			Msg("Could not mark 'agent' flag required")
		os.Exit(1)
	}
}

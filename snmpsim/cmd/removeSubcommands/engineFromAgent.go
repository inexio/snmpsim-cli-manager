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
		baseUrl := viper.GetString("mgmt.http.baseUrl")
		username := viper.GetString("mgmt.http.authUsername")
		password := viper.GetString("mgmt.http.authPassword")

		//Create a new client
		client, err := snmpsimclient.NewManagementClient(baseUrl)
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

		//Read in the engine-id
		engineId, err := cmd.Flags().GetInt("engine")
		if err != nil {
			log.Error().
				Msg("Error while retrieving engineId")
			os.Exit(1)
		}

		//Read in the agent-id
		agentId, err := cmd.Flags().GetInt("agent")
		if err != nil {
			log.Error().
				Msg("Error while retrieving agentId")
			os.Exit(1)
		}

		//Remove the engine from the agent
		err = client.RemoveEngineFromAgent(engineId, agentId)
		if err != nil {
			log.Error().
				Msg("Error while removing the engine from the agent")
			os.Exit(1)
		}
		fmt.Println("Engine", engineId, "has been removed from agent", agentId)
	},
}

func init() {
	//Set engine flag
	EngineFromAgentCmd.Flags().Int("engine", 0, "Id of the engine the will be remove from the agent")
	err := EngineFromAgentCmd.MarkFlagRequired("engine")
	if err != nil {
		log.Error().
			Msg("Could not mark 'engine' flag required")
		os.Exit(1)
	}

	//Set agent flag
	EngineFromAgentCmd.Flags().Int("agent", 0, "Id of the agent the engine will be removed from")
	err = EngineFromAgentCmd.MarkFlagRequired("agent")
	if err != nil {
		log.Error().
			Msg("Could not mark 'agent' flag required")
		os.Exit(1)
	}
}

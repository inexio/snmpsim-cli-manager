package addsubcommands

import (
	"fmt"
	"os"

	snmpsimclient "github.com/inexio/snmpsim-restapi-go-client"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// EngineToAgentCmd represents the agentToAgent command
var EngineToAgentCmd = &cobra.Command{
	Use:   "engine-to-agent",
	Args:  cobra.ExactArgs(0),
	Short: "Adds an engine to an agent",
	Long:  `Adds the engine with a given engine-id to the agent with the given agent-id.`,
	Run: func(cmd *cobra.Command, args []string) {
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
		if username != "" && password != "" {
			err = client.SetUsernameAndPassword(username, password)
			if err != nil {
				log.Error().
					Msg("Error while setting username and password")
				os.Exit(1)
			}
		}

		//Read in the engine-id
		engineID, err := cmd.Flags().GetInt("engine")
		if err != nil {
			log.Error().
				Msg("Error while retrieving engine-id")
			os.Exit(1)
		}

		//Read in the agent-id
		agentID, err := cmd.Flags().GetInt("agent")
		if err != nil {
			log.Error().
				Msg("Error while retrieving agent-id")
			os.Exit(1)
		}

		//Add the engine to the agent
		err = client.AddEngineToAgent(agentID, engineID)
		if err != nil {
			log.Error().
				Msg("Error while adding engine to agent")
			os.Exit(1)
		}
		fmt.Println("Engine", engineID, "has been added to agent ", agentID)
	},
}

func init() {
	//Set engine flag
	EngineToAgentCmd.Flags().Int("engine", 0, "Id of the engine that is to be added to the agent")
	err := EngineToAgentCmd.MarkFlagRequired("engine")
	if err != nil {
		log.Error().
			Msg("Could not mark 'engine' flag required")
		os.Exit(1)
	}

	//Set agent flag
	EngineToAgentCmd.Flags().Int("agent", 0, "Id of the agent to that the engine will be added")
	err = EngineToAgentCmd.MarkFlagRequired("agent")
	if err != nil {
		log.Error().
			Msg("Could not mark 'agent' flag required")
		os.Exit(1)
	}
}

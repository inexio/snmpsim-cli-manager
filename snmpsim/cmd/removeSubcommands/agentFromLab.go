package removesubcommands

import (
	"fmt"
	"os"

	snmpsimclient "github.com/inexio/snmpsim-restapi-go-client"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// AgentFromLabCmd represents the agentFromLab command
var AgentFromLabCmd = &cobra.Command{
	Use:   "agent-from-lab",
	Args:  cobra.ExactArgs(0),
	Short: "Removes an agent from a lab",
	Long:  `Removes the agent with the given agent-id from the lab with the given lab-id`,
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
		err = client.SetUsernameAndPassword(username, password)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error while setting username and password")
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

		//Read in the lab-id
		labID, err := cmd.Flags().GetInt("lab")
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error while retrieving labID")
			os.Exit(1)
		}

		//Remove the agent from the lab
		err = client.RemoveAgentFromLab(labID, agentID)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error while removing the agent from the lab")
			os.Exit(1)
		}
		fmt.Println("Agent", agentID, "has been removed from lab", labID)
	},
}

func init() {
	//Set agent flag
	AgentFromLabCmd.Flags().Int("agent", 0, "Id of the agent that is to be removed from the lab")
	err := AgentFromLabCmd.MarkFlagRequired("agent")
	if err != nil {
		log.Error().
			Err(err).
			Msg("Could not mark 'agent' flag required")
		os.Exit(1)
	}

	//Set lab flag
	AgentFromLabCmd.Flags().Int("lab", 0, "Id of the lab the agent will be removed from")
	err = AgentFromLabCmd.MarkFlagRequired("lab")
	if err != nil {
		log.Error().
			Err(err).
			Msg("Could not mark 'lab' flag required")
		os.Exit(1)
	}
}

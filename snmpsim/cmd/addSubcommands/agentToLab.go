package addsubcommands

import (
	"fmt"
	"os"

	snmpsimclient "github.com/inexio/snmpsim-restapi-go-client"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// AgentToLabCmd represents the agentToLab command
var AgentToLabCmd = &cobra.Command{
	Use:   "agent-to-lab",
	Args:  cobra.ExactArgs(0),
	Short: "Adds an agent to a lab",
	Long:  `Adds the agent with the given agent-id to the lab with the given lab-id.`,
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
		err = client.SetUsernameAndPassword(username, password)
		if err != nil {
			log.Error().
				Msg("Error while setting username and password")
			os.Exit(1)
		}

		//Read in the agent-id
		agentID, err := cmd.Flags().GetInt("agent")
		if err != nil {
			log.Error().
				Msg("Error while retrieving agent-id")
			os.Exit(1)
		}

		//Read in the lab-id
		labID, err := cmd.Flags().GetInt("lab")
		if err != nil {
			log.Error().
				Msg("Error while retrieving lab-id")
			os.Exit(1)
		}

		//Add the agent to the lab
		err = client.AddAgentToLab(labID, agentID)
		if err != nil {
			log.Error().
				Msg("Error while adding agent to lab")
			os.Exit(1)
		}
		fmt.Println("Agent", agentID, "has been added to lab", labID)
	},
}

func init() {
	//Set agent flag
	AgentToLabCmd.Flags().Int("agent", 0, "Id of the agent that is to be added to the lab")
	err := AgentToLabCmd.MarkFlagRequired("agent")
	if err != nil {
		log.Error().
			Msg("Could not mark 'agent' flag required")
		os.Exit(1)
	}

	//Set lab flag
	AgentToLabCmd.Flags().Int("lab", 0, "Id of the lab to that the agent will be added")
	err = AgentToLabCmd.MarkFlagRequired("lab")
	if err != nil {
		log.Error().
			Msg("Could not mark 'lab' flag required")
		os.Exit(1)
	}
}

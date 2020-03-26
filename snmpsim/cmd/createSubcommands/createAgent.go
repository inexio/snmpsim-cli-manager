package createsubcommands

import (
	"fmt"
	"github.com/inexio/snmpsim-restapi-go-client"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// CreateAgentCmd represents the createAgent command
var CreateAgentCmd = &cobra.Command{
	Use:   "agent",
	Args:  cobra.ExactArgs(0),
	Short: "Creates a new agent",
	Long:  `Creates a new agent and returns its id.`,
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

		//Read in the agents name and data-directory
		name := cmd.Flag("name").Value.String()
		dataDir := cmd.Flag("dataDir").Value.String()

		//Create an agent
		var agent snmpsimclient.Agent
		if cmd.Flag("tag").Changed {
			//Read in the tag-id
			tagId, err := cmd.Flags().GetInt("tag")
			if err != nil {
				log.Error().
					Msg("Error while retrieving tagId")
				os.Exit(1)
			}

			//Validate tag-id
			if tagId == 0 {
				log.Error().
					Msg("TagId can not be 0")
				os.Exit(1)
			}

			//Check if tag with given id exists
			_, err = client.GetTag(tagId)
			if err != nil {
				log.Error().
					Msg("No tag with the given id found")
				os.Exit(1)
			}

			agent, err = client.CreateAgentWithTag(name, dataDir, tagId)
			if err != nil {
				log.Error().
					Msg("Error during creation of the agent")
				os.Exit(1)
			}
		} else {
			agent, err = client.CreateAgent(name, dataDir)
			if err != nil {
				log.Error().
					Msg("Error during creation of the agent")
				os.Exit(1)
			}
		}

		fmt.Println("Successfully created agent.")
		fmt.Println("Id:", agent.Id)

		//Add agent to lab (if lab flag is set)
		if cmd.Flag("lab").Changed {
			//Read in lab-id
			labId, err := cmd.Flags().GetInt("lab")
			if err != nil {
				log.Error().
					Msg("Error while retrieving lab-id")
				os.Exit(1)
			}

			//Check if lab with given id exists
			_, err = client.GetLab(labId)
			if err != nil {
				log.Error().
					Msg("No lab with the given id found")
				os.Exit(1)
			}

			//Add agent to lab
			err = client.AddAgentToLab(labId, agent.Id)
			if err != nil {
				log.Error().
					Msg("Error while adding agent to lab")
				os.Exit(1)
			}
			fmt.Println("Successfully added agent", agent.Id, "to lab", labId)
		}
	},
}

func init() {
	CreateAgentCmd.Flags().String("dataDir", "", "Data directory of the agent")
	err := CreateAgentCmd.MarkFlagRequired("dataDir")
	if err != nil {
		log.Error().
			Msg("Could not mark 'dataDir' flag required")
		os.Exit(1)
	}
	CreateAgentCmd.Flags().Int("tag", 0, "Optional flag to mark an agent with a tag (requires a tag-id)")
	CreateAgentCmd.Flags().Int("lab", 0, "Optional flag to add the agent to a lab (requires a lab-id)")
}

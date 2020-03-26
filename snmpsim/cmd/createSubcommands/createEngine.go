package createsubcommands

import (
	"fmt"
	snmpsimclient "github.com/inexio/snmpsim-restapi-go-client"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// CreateEngineCmd represents the createEngine command
var CreateEngineCmd = &cobra.Command{
	Use:   "engine",
	Args:  cobra.ExactArgs(0),
	Short: "Creates a new engine",
	Long:  `Creates a new engine and returns its id`,
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

		//Read in the engines name and engineId
		name := cmd.Flag("name").Value.String()
		engineId := cmd.Flag("engineId").Value.String()

		//Create an engine
		var engine snmpsimclient.Engine
		if cmd.Flag("tag").Changed {
			//Read in tag-id
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

			engine, err = client.CreateEngineWithTag(name, engineId, tagId)
			if err != nil {
				log.Error().
					Msg("Error during creation of the engine")
				os.Exit(1)
			}
		} else {
			engine, err = client.CreateEngine(name, engineId)
			if err != nil {
				log.Error().
					Msg("Error during creation of the engine")
				os.Exit(1)
			}
		}

		fmt.Println("Engine has been created successfully.")
		fmt.Println("Id:", engine.Id)

		//Add engine to agent (if agent flag is set)
		if cmd.Flag("agent").Changed {
			//Read in agent-id
			agentId, err := cmd.Flags().GetInt("agent")
			if err != nil {
				log.Error().
					Msg("Error while retrieving agentId")
				os.Exit(1)
			}

			//Check if agent with given id exists
			_, err = client.GetAgent(agentId)
			if err != nil {
				log.Error().
					Msg("No agent with the given id found")
				os.Exit(1)
			}

			//Add engine to agent
			err = client.AddEngineToAgent(agentId, engine.Id)
			if err != nil {
				log.Error().
					Msg("Error while adding engine to agent")
				os.Exit(1)
			}
			fmt.Println("Successfully added engine", engine.Id, "to agent ", agentId)
		}
	},
}

func init() {
	CreateEngineCmd.Flags().String("engineId", "", "Freely selectable engine-id (not the internal id)")
	err := CreateEngineCmd.MarkFlagRequired("engineId")
	if err != nil {
		log.Error().
			Msg("Could not mark 'engineId' flag required")
		os.Exit(1)
	}

	CreateEngineCmd.Flags().Int("tag", 0, "Optional flag to mark an engine with a tag (requires a tag-id)")
	CreateEngineCmd.Flags().Int("agent", 0, "Optional flag to add the engine to an agent (requires an agent-id)")
}

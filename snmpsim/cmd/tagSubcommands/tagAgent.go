package tagsubcommands

import (
	"fmt"
	snmpsimclient "github.com/inexio/snmpsim-restapi-go-client"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strconv"
)

// TagAgentCmd represents the tagAgent command
var TagAgentCmd = &cobra.Command{
	Use:   "agent <agent-id> --tag <tag-id>",
	Args:  cobra.ExactArgs(1),
	Short: "Tags an agent",
	Long:  `Tags the agent with the given agent-id with the given tag-id.`,
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

		//Read in the tag-id
		tagID, err := cmd.Flags().GetInt("tag")
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error while converting tagID from string to integer")
			os.Exit(1)
		}

		//Read in the agent-id
		agentID, err := strconv.Atoi(args[0])
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error while converting agentID from string to integer")
			os.Exit(1)
		}

		//Add the tag to the agent
		err = client.AddTagToAgent(agentID, tagID)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error during addition of the the tag to the agent")
			os.Exit(1)
		}
		fmt.Println("Agent", agentID, "has been tagged with tag", tagID)
	},
}

package untagsubcommands

import (
	"fmt"
	snmpsimclient "github.com/inexio/snmpsim-restapi-go-client"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strconv"
)

// UntagEngineCmd represents the untagEngine command
var UntagEngineCmd = &cobra.Command{
	Use:   "engine <engine-id> --tag <tag-id>",
	Args:  cobra.ExactArgs(1),
	Short: "Untags an engine",
	Long:  `Removes the tag with the given tag-id from the engine with the given engine-id.`,
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

		//Read in the tag-id
		tagID, err := cmd.Flags().GetInt("tag")
		if err != nil {
			log.Error().
				Msg("Error while converting tagID from string to integer")
			os.Exit(1)
		}

		//Read in the engine-id
		engineID, err := strconv.Atoi(args[0])
		if err != nil {
			log.Error().
				Msg("Error while converting engine from string to integer")
			os.Exit(1)
		}

		//Removing the tag from the engine
		err = client.RemoveTagFromEngine(engineID, tagID)
		if err != nil {
			log.Error().
				Msg("Error while removing tag from the engine")
			os.Exit(1)
		}
		fmt.Println("Tag", tagID, "has been removed from engine", engineID)
	},
}

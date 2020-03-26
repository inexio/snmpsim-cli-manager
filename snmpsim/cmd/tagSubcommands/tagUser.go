package tagsubcommands

import (
	"fmt"
	snmpsimclient "github.com/inexio/snmpsim-restapi-go-client"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

// TagUserCmd represents the tagUser command
var TagUserCmd = &cobra.Command{
	Use:   "user <user-id> --tag <tag-id>",
	Args:  cobra.ExactArgs(1),
	Short: "Tags an user",
	Long:  `Tags the user with the given user-id with the given tag-id.`,
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

		//Read in the tag-id
		tagId, err := cmd.Flags().GetInt("tag")
		if err != nil {
			log.Error().
				Msg("Error while converting tagId from string to integer")
			os.Exit(1)
		}

		//Read in the user-id
		userId, err := strconv.Atoi(args[0])
		if err != nil {
			log.Error().
				Msg("Error while converting userId from string to integer")
			os.Exit(1)
		}

		//Add the tag to the user
		err = client.AddTagToUser(userId, tagId)
		if err != nil {
			log.Error().
				Msg("Error while adding tag to the user")
			os.Exit(1)
		}
		fmt.Println("User", userId, "has been added tagged with tag", tagId)
	},
}

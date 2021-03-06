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

// UntagUserCmd represents the untagUser command
var UntagUserCmd = &cobra.Command{
	Use:   "user <user-id> --tag <tag-id>",
	Args:  cobra.ExactArgs(1),
	Short: "Untags an user",
	Long:  `Removes the tag with the given tag-id from the user with the given user-id.`,
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

		//Read in the user-id
		userID, err := strconv.Atoi(args[0])
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error while converting userID from string to integer")
			os.Exit(1)
		}

		//Removing the tag from the user
		err = client.RemoveTagFromUser(userID, tagID)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error while removing tag from the user")
			os.Exit(1)
		}
		fmt.Println("Tag", tagID, "has been removed from user", userID)
	},
}

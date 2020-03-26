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

// UntagLabCmd represents the untagLab command
var UntagLabCmd = &cobra.Command{
	Use:   "lab <lab-id> --tag <tag-id>",
	Args:  cobra.ExactArgs(1),
	Short: "Untags a lab",
	Long:  `Removes the tag with the given tag-id from the lab with the given lab-id.`,
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

		//Read in the lab-id
		labId, err := strconv.Atoi(args[0])
		if err != nil {
			log.Error().
				Msg("Error while converting labId from string to integer")
			os.Exit(1)
		}

		//Removing the tag from tha lab
		err = client.RemoveTagFromLab(labId, tagId)
		if err != nil {
			log.Error().
				Msg("Error while removing tag from the lab")
			os.Exit(1)
		}
		fmt.Println("Tag", tagId, "has been removed from lab", labId)
	},
}

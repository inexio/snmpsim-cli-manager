package deletesubcommands

import (
	"fmt"
	snmpsimclient "github.com/inexio/snmpsim-restapi-go-client"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strconv"
)

// DeleteTagCmd represents the deleteTag command
var DeleteTagCmd = &cobra.Command{
	Use:   "tag <id>",
	Args:  cobra.ExactArgs(1),
	Short: "Deletes a tag",
	Long:  `Deletes the tag with the given tag-id`,
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

		//Read in the tag-id
		tagID, err := strconv.Atoi(args[0])
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error during conversion of tagID from string to integer")
			os.Exit(1)
		}

		//Delete the tag
		err = client.DeleteTag(tagID)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error while deleting tag")
			os.Exit(1)
		}

		fmt.Println("Tag", args[0], "has been deleted successfully.")
	},
}

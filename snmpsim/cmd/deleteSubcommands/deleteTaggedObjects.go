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

// DeleteTaggedObjectsCmd represents the deleteTaggedObjects command
var DeleteTaggedObjectsCmd = &cobra.Command{
	Use:   "tagged-objects <tag_id>",
	Args:  cobra.ExactArgs(1),
	Short: "Deletes all tagged objects",
	Long:  `Deletes all objects tagged with the given tag-id`,
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
		tagID, err := strconv.Atoi(args[0])
		if err != nil {
			log.Error().
				Msg("Error during conversion of tagID from string to integer")
			os.Exit(1)
		}

		//Delete all objects tagged with given id
		tag, err := client.DeleteAllObjectsWithTag(tagID)
		if err != nil {
			log.Error().
				Msg("Error during deletion of all tagged objects")
			os.Exit(1)
		}

		fmt.Println("All Objects tagged with", tag.Name, "(Id:", tag.ID, ") were deleted successfully.")
	},
}

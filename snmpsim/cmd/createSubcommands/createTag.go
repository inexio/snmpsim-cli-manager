package createsubcommands

import (
	"fmt"
	snmpsimclient "github.com/inexio/snmpsim-restapi-go-client"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// CreateTagCmd represents the createTag command
var CreateTagCmd = &cobra.Command{
	Use:   "tag",
	Args:  cobra.ExactArgs(0),
	Short: "Creates a new tag",
	Long:  `Creates a new tag and returns its id`,
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

		//Read in the tags name and description
		name := cmd.Flag("name").Value.String()
		description := cmd.Flag("description").Value.String()

		//Create the tag
		var tag snmpsimclient.Tag
		tag, err = client.CreateTag(name, description)
		if err != nil {
			log.Error().
				Msg("Error during creation of the tag")
			os.Exit(1)
		}

		fmt.Println("Tag has been created successfully.")
		fmt.Println("Id:", tag.Id)
	},
}

func init() {
	CreateTagCmd.Flags().String("description", "", "Description of the tag")
	err := CreateTagCmd.MarkFlagRequired("description")
	if err != nil {
		log.Error().
			Msg("Could not mark 'description' flag required")
		os.Exit(1)
	}
}

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

// TagLabCmd represents the tagLab command
var TagLabCmd = &cobra.Command{
	Use:   "lab <lab-id> --tag <tag-id>",
	Args:  cobra.ExactArgs(1),
	Short: "Tags a lab",
	Long:  `Tags the lab with the given lab-id with the given tag-id.`,
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

		//Read in the lab-id
		labID, err := strconv.Atoi(args[0])
		if err != nil {
			log.Error().
				Msg("Error while converting labID from string to integer")
			os.Exit(1)
		}

		//Add the tag to the lab
		err = client.AddTagToLab(labID, tagID)
		if err != nil {
			log.Error().
				Msg("Error while adding tag to the lab")
			os.Exit(1)
		}
		fmt.Println("Lab", labID, "has been tagged with tag", tagID)
	},
}

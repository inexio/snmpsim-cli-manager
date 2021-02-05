package createsubcommands

import (
	"fmt"
	snmpsimclient "github.com/inexio/snmpsim-restapi-go-client"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// CreateLabCmd represents the createLab command
var CreateLabCmd = &cobra.Command{
	Use:   "lab",
	Args:  cobra.ExactArgs(0),
	Short: "Creates a new lab",
	Long:  `Creates a new lab and returns its id`,
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

		//Read in the labs name
		name := cmd.Flag("name").Value.String()

		//Create a lab
		var lab snmpsimclient.Lab
		if cmd.Flag("tag").Changed {
			//Read in tag-id
			tagID, err := cmd.Flags().GetInt("tag")
			if err != nil {
				log.Error().
					Err(err).
					Msg("Error while retrieving tagID")
				os.Exit(1)
			}

			//Validate tag-id
			if tagID == 0 {
				log.Error().
					Msg("tagID can not be 0")
				os.Exit(1)
			}

			//Check if tag with given id exists
			_, err = client.GetTag(tagID)
			if err != nil {
				log.Error().
					Err(err).
					Msg("No tag with the given id found")
				os.Exit(1)
			}

			lab, err = client.CreateLabWithTag(name, tagID)
			if err != nil {
				log.Error().
					Err(err).
					Msg("Error during creation of the lab")
				os.Exit(1)
			}
		} else {
			lab, err = client.CreateLab(name)
			if err != nil {
				log.Error().
					Err(err).
					Msg("Error during creation of the lab")
				os.Exit(1)
			}
		}

		fmt.Println("Lab has been created successfully.")
		fmt.Println("Id:", lab.ID)
	},
}

func init() {
	CreateLabCmd.Flags().Int("tag", 0, "Optional flag to mark a lab with a tag (requires a tag-id)")
}

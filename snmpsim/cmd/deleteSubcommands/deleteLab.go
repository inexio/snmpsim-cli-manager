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

// DeleteLabCmd represents the deleteLab command
var DeleteLabCmd = &cobra.Command{
	Use:   "lab <id>",
	Args:  cobra.ExactArgs(1),
	Short: "Deletes a lab",
	Long:  `Deletes the lab with the given lab-id`,
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

		//Read in the lab-id
		labID, err := strconv.Atoi(args[0])
		if err != nil {
			log.Error().
				Msg("Error during conversion of labID from string to integer")
			os.Exit(1)
		}

		//Delete the lab
		err = client.DeleteLab(labID)
		if err != nil {
			log.Error().
				Msg("Error while deleting lab")
			os.Exit(1)
		}

		fmt.Println("Lab", args[0], "has been deleted successfully.")
	},
}

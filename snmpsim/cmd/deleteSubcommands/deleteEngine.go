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

// DeleteEngineCmd represents the deleteEngine command
var DeleteEngineCmd = &cobra.Command{
	Use:   "engine <id>",
	Args:  cobra.ExactArgs(1),
	Short: "Deletes an engine",
	Long:  `Deletes the engine with the given engine-id`,
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

		//Read in the engine-id
		engineID, err := strconv.Atoi(args[0])
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error during conversion of engineID from string to integer")
			os.Exit(1)
		}

		//Delete the engine
		err = client.DeleteEngine(engineID)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error while deleting engine")
			os.Exit(1)
		}

		fmt.Println("Engine", args[0], "has been deleted successfully.")
	},
}

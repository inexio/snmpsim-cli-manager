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

// DeleteEndpointCmd represents the deleteEndpoint command
var DeleteEndpointCmd = &cobra.Command{
	Use:   "endpoint <id>",
	Args:  cobra.ExactArgs(1),
	Short: "Deletes an endpoint",
	Long:  `Deletes the endpoint with the given endpoint-id`,
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

		//Read in the endpoint-id
		endpointId, err := strconv.Atoi(args[0])
		if err != nil {
			log.Error().
				Msg("Error during conversion of tagId from string to integer")
			os.Exit(1)
		}

		//Delete the endpoint
		err = client.DeleteEndpoint(endpointId)
		if err != nil {
			log.Error().
				Msg("Error while deleting endpoint")
			os.Exit(1)
		}

		fmt.Println("Endpoint", args[0], "has been deleted successfully.")
	},
}

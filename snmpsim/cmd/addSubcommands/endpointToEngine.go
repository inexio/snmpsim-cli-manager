package addsubcommands

import (
	"fmt"
	"os"

	snmpsimclient "github.com/inexio/snmpsim-restapi-go-client"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// EndpointToEngineCmd represents the endpointToEngine command
var EndpointToEngineCmd = &cobra.Command{
	Use:   "endpoint-to-engine",
	Args:  cobra.ExactArgs(0),
	Short: "Adds an endpoint to an engine",
	Long:  `Adds the endpoint with the given endpoint-id to the engine with the given engine-id.`,
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

		//Read in the endpoint-id
		endpointID, err := cmd.Flags().GetInt("endpoint")
		if err != nil {
			log.Error().
				Msg("Error while retrieving endpoint-id")
			os.Exit(1)
		}

		//Read in the engine-id
		engineID, err := cmd.Flags().GetInt("engine")
		if err != nil {
			log.Error().
				Msg("Error while retrieving engine-id")
			os.Exit(1)
		}

		//Add the endpoint to the engine
		err = client.AddEndpointToEngine(engineID, endpointID)
		if err != nil {
			log.Error().
				Msg("Error while adding endpoint to engine")
			os.Exit(1)
		}
		fmt.Println("Endpoint", endpointID, "has been added to engine", engineID)
	},
}

func init() {
	//Set endpoint flag
	EndpointToEngineCmd.Flags().Int("endpoint", 0, "Id of the endpoint that is to be added to the engine")
	err := EndpointToEngineCmd.MarkFlagRequired("endpoint")
	if err != nil {
		log.Error().
			Msg("Could not mark 'endpoint' flag required")
		os.Exit(1)
	}

	//Set engine flag
	EndpointToEngineCmd.Flags().Int("engine", 0, "Id of the engine to that the endpoint will be added")
	err = EndpointToEngineCmd.MarkFlagRequired("engine")
	if err != nil {
		log.Error().
			Msg("Could not mark 'engine' flag required")
		os.Exit(1)
	}
}

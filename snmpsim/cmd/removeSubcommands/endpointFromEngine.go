package removesubcommands

import (
	"fmt"
	"os"

	snmpsimclient "github.com/inexio/snmpsim-restapi-go-client"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// EndpointFromEngineCmd represents the endpointFromEngine command
var EndpointFromEngineCmd = &cobra.Command{
	Use:   "endpoint-from-engine",
	Args:  cobra.ExactArgs(0),
	Short: "Removes an endpoint from an engine",
	Long:  `Removes the endpoint with the given endpoint-id from the engine with the given engine-id`,
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

		//Read in the endpoint-id
		endpointID, err := cmd.Flags().GetInt("endpoint")
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error while retrieving endpointID")
			os.Exit(1)
		}

		//Read in the engine-id
		engineID, err := cmd.Flags().GetInt("engine")
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error while retrieving engineID")
			os.Exit(1)
		}

		//Remove the endpoint from the engine
		err = client.RemoveEndpointFromEngine(engineID, endpointID)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error while removing the endpoint from the engine")
			os.Exit(1)
		}
		fmt.Println("Endpoint", endpointID, "has been removed from engine", engineID)
	},
}

func init() {
	//Set endpoint flag
	EndpointFromEngineCmd.Flags().Int("endpoint", 0, "Id of the endpoint that is to be removed from the engine")
	err := EndpointFromEngineCmd.MarkFlagRequired("endpoint")
	if err != nil {
		log.Error().
			Err(err).
			Msg("Could not mark 'endpoint' flag required")
		os.Exit(1)
	}

	//Set engine flag
	EndpointFromEngineCmd.Flags().Int("engine", 0, "Id of the engine from that the endpoint will be removed")
	err = EndpointFromEngineCmd.MarkFlagRequired("engine")
	if err != nil {
		log.Error().
			Err(err).
			Msg("Could not mark 'engine' flag required")
		os.Exit(1)
	}
}

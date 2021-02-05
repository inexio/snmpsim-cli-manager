package createsubcommands

import (
	"fmt"
	snmpsimclient "github.com/inexio/snmpsim-restapi-go-client"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// CreateEndpointCmd represents the createEndpoint command
var CreateEndpointCmd = &cobra.Command{
	Use:   "endpoint",
	Args:  cobra.ExactArgs(0),
	Short: "Creates a new endpoint",
	Long:  `Creates a new endpoint and returns its id.`,
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

		//Read in the endpoints name, address and protocol
		name := cmd.Flag("name").Value.String()
		address := cmd.Flag("address").Value.String()
		protocol := cmd.Flag("protocol").Value.String()

		//Create an endpoint
		var endpoint snmpsimclient.Endpoint
		if cmd.Flag("tag").Changed {
			//Read in the tag-id
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

			endpoint, err = client.CreateEndpointWithTag(name, address, protocol, tagID)
			if err != nil {
				log.Error().
					Err(err).
					Msg("Error during creation of the endpoint")
				os.Exit(1)
			}
		} else {
			endpoint, err = client.CreateEndpoint(name, address, protocol)
			if err != nil {
				log.Error().
					Err(err).
					Msg("Error during creation of the endpoint")
				os.Exit(1)
			}
		}

		fmt.Println("Endpoint has been created successfully.")
		fmt.Println("Id:", endpoint.ID)

		//Add endpoint to engine (if engine flag is set)
		if cmd.Flag("engine").Changed {
			//Read in engine-id
			engineID, err := cmd.Flags().GetInt("engine")
			if err != nil {
				log.Error().
					Err(err).
					Msg("Error while retrieving engine-id")
				os.Exit(1)
			}

			//Check if engine with given id exists
			_, err = client.GetEngine(engineID)
			if err != nil {
				log.Error().
					Err(err).
					Msg("No engine with the given id found")
				os.Exit(1)
			}

			//Add endpoint to engine
			err = client.AddEndpointToEngine(engineID, endpoint.ID)
			if err != nil {
				log.Error().
					Err(err).
					Msg("Error while adding endpoint to engine")
				os.Exit(1)
			}
			fmt.Println("Successfully added endpoint", endpoint.ID, "to engine", engineID)
		}
	},
}

func init() {
	CreateEndpointCmd.Flags().String("address", "", "The address of the endpoint")
	err := CreateEndpointCmd.MarkFlagRequired("address")
	if err != nil {
		log.Error().
			Err(err).
			Msg("Could not mark 'address' flag required")
		os.Exit(1)
	}

	CreateEndpointCmd.Flags().String("protocol", "", "The protocol of the endpoint")
	err = CreateEndpointCmd.MarkFlagRequired("protocol")
	if err != nil {
		log.Error().
			Err(err).
			Msg("Could not mark 'protocol' flag required")
		os.Exit(1)
	}

	CreateEndpointCmd.Flags().Int("tag", 0, "Optional flag to mark an endpoint with a tag (requires a tag-id)")
	CreateEndpointCmd.Flags().Int("engine", 0, "Optional flag to add the endpoint to an engine (requires an engine-id)")
}

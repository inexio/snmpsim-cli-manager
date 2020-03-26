package getsubcommands

import (
	snmpsimclient "github.com/inexio/snmpsim-restapi-go-client"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strconv"
)

// GetEndpointCmd represents the getEndpoint command
var GetEndpointCmd = &cobra.Command{
	Use:   "endpoint <id>",
	Args:  cobra.ExactArgs(1),
	Short: "Returns details about an endpoint",
	Long: `Returns details about the endpoint with the given id.

These details include:
	- id
	- name
	- protocol
	- address
	- tags`,
	Run: func(cmd *cobra.Command, args []string) {
		//Parse all persistent flags
		format, depth, prettified := parsePersistentFlags(cmd)

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

		//Read in and convert the id
		id, err := strconv.Atoi(args[0])
		if err != nil {
			log.Error().
				Msg("Error during conversion of id from string to integer")
			os.Exit(1)
		}

		//Get and print the endpoint
		var endpoint snmpsimclient.Endpoint
		endpoint, err = client.GetEndpoint(id)
		if err != nil {
			log.Error().
				Msg("Error while getting endpoint")
			os.Exit(1)
		}
		err = printData(endpoint, format, prettified, depth)
		if err != nil {
			log.Error().
				Msg("Error while printing data")
			os.Exit(1)
		}
	},
}

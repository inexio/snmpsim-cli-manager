package getsubcommands

import (
	snmpsimclient "github.com/inexio/snmpsim-restapi-go-client"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// GetEndpointsCmd represents the getEndpoints command
var GetEndpointsCmd = &cobra.Command{
	Use:   "endpoints",
	Args:  cobra.ExactArgs(0),
	Short: "Returns a list of all endpoints",
	Long: `Returns a detailed list of all endpoints available.

All details of one specific endpoint can be retrieved via 'get endpoint <endpoint-id>'.`,
	Run: func(cmd *cobra.Command, args []string) {
		//Parse all persistent flags
		format, depth, prettified := parsePersistentFlags(cmd)

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

		//Parse filters from flags
		filters := parseFilters(cmd)

		//getting and printing the endpoint
		var endpoints snmpsimclient.Endpoints
		endpoints, err = client.GetEndpoints(filters)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error while getting agents")
			os.Exit(1)
		}
		err = printData(endpoints, format, prettified, depth)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error while printing data")
			os.Exit(1)
		}
	},
}

func init() {
	GetEndpointsCmd.Flags().String("name", "", "Set the name filter")
	GetEndpointsCmd.Flags().String("protocol", "", "Set the protocol filter")
	GetEndpointsCmd.Flags().String("address", "", "Set the address filter")
}

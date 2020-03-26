package getsubcommands

import (
	snmpsimclient "github.com/inexio/snmpsim-restapi-go-client"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// GetEnginesCmd represents the getEngines command
var GetEnginesCmd = &cobra.Command{
	Use:   "engines",
	Args:  cobra.ExactArgs(0),
	Short: "Returns a list of all engines",
	Long: `Returns a detailed list of all engines available.

All details of one specific engine can be retrieved via 'get engine <engine-id>'.`,
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

		//Parse filters from flags
		filters := parseFilters(cmd)

		//Get and print the engines
		var engines snmpsimclient.Engines
		engines, err = client.GetEngines(filters)
		if err != nil {
			log.Error().
				Msg("Error while getting engines")
			os.Exit(1)
		}
		err = printData(engines, format, prettified, depth)
		if err != nil {
			log.Error().
				Msg("Error while printing data")
			os.Exit(1)
		}
	},
}

func init() {
	GetEnginesCmd.Flags().String("name", "", "Set the name filter")
	GetEnginesCmd.Flags().String("engine_id", "", "Set the engine id filter")
}

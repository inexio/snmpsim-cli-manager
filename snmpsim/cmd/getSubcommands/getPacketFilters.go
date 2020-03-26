package getsubcommands

import (
	snmpsimclient "github.com/inexio/snmpsim-restapi-go-client"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// PacketFiltersCmd represents the packetFilters command
var PacketFiltersCmd = &cobra.Command{
	Use:   "packet-filters",
	Args:  cobra.ExactArgs(0),
	Short: "Returns all possible packet filters",
	Long:  `Returns a complete list of all possible filters for the packets endpoint`,
	Run: func(cmd *cobra.Command, args []string) {
		//Parse all persistent flags
		format, depth, prettified := parsePersistentFlags(cmd)

		//Load the client data from the config
		baseUrl := viper.GetString("metrics.http.baseUrl")
		username := viper.GetString("metrics.http.authUsername")
		password := viper.GetString("metrics.http.authPassword")

		//Create a new client
		client, err := snmpsimclient.NewMetricsClient(baseUrl)
		if err != nil {
			log.Error().
				Msg("Error during creation of new metrics client")
			os.Exit(1)
		}
		err = client.SetUsernameAndPassword(username, password)
		if err != nil {
			log.Error().
				Msg("Error while setting username and password")
			os.Exit(1)
		}

		//getting and printing packet-filters
		var filters snmpsimclient.PacketFilters
		filters, err = client.GetPacketFilters()
		if err != nil {
			log.Error().
				Msg("Error while getting packet-filters")
			os.Exit(1)
		}
		err = printData(filters, format, prettified, depth)
		if err != nil {
			log.Error().
				Msg("Error while printing data")
			os.Exit(1)
		}
	},
}

package getsubcommands

import (
	snmpsimclient "github.com/inexio/snmpsim-restapi-go-client"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// PacketsCmd represents the packets command
var PacketsCmd = &cobra.Command{
	Use:   "packets",
	Args:  cobra.ExactArgs(0),
	Short: "Returns all network packets",
	Long:  `Returns a list of all network transport packets including successful and unsuccessful snmp commands.`,
	Run: func(cmd *cobra.Command, args []string) {
		//Parse all persistent flags
		format, depth, prettified := parsePersistentFlags(cmd)

		//Load the client data from the config
		baseURL := viper.GetString("metrics.http.baseURL")
		username := viper.GetString("metrics.http.authUsername")
		password := viper.GetString("metrics.http.authPassword")

		//Create a new client
		client, err := snmpsimclient.NewMetricsClient(baseURL)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error during creation of new metrics client")
			os.Exit(1)
		}
		err = client.SetUsernameAndPassword(username, password)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error while setting username and password")
			os.Exit(1)
		}

		//parse the filters from the commands flags
		filters := parseFilters(cmd)

		//Get and print the packets
		var packets snmpsimclient.PacketMetrics
		packets, err = client.GetPackets(filters)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error while getting packets")
			os.Exit(1)
		}
		err = printData(packets, format, prettified, depth)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error while printing data")
			os.Exit(1)
		}
	},
}

func init() {
	PacketsCmd.Flags().String("protocol", "", "Setting the protocol message filter")
	PacketsCmd.Flags().String("local_address", "", "Setting the local_address message filter")
	PacketsCmd.Flags().String("peer_address", "", "Setting the peer_address message filter")
}

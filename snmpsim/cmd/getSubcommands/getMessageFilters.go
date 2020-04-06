package getsubcommands

import (
	snmpsimclient "github.com/inexio/snmpsim-restapi-go-client"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	//"strconv"
)

// MessageFiltersCmd represents the messageFilters command
var MessageFiltersCmd = &cobra.Command{
	Use:   "message-filters",
	Args:  cobra.ExactArgs(0),
	Short: "Returns all possible message filters",
	Long:  `Returns a complete list of all possible filters for the messages endpoint`,
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
				Msg("Error during creation of new metrics client")
			os.Exit(1)
		}
		err = client.SetUsernameAndPassword(username, password)
		if err != nil {
			log.Error().
				Msg("Error while setting username and password")
			os.Exit(1)
		}

		//Get and print the message-filters
		var filters snmpsimclient.MessageFilters
		filters, err = client.GetMessageFilters()
		if err != nil {
			log.Error().
				Msg("Error while getting message-filters")
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

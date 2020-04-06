package getsubcommands

import (
	snmpsimclient "github.com/inexio/snmpsim-restapi-go-client"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// MessagesCmd represents the messages command
var MessagesCmd = &cobra.Command{
	Use:   "messages",
	Args:  cobra.ExactArgs(0),
	Short: "Returns all network messages",
	Long:  `Returns a list of all network transport messages including only successful snmp commands.`,
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

		//Parse filters from flags
		filters := parseFilters(cmd)

		//Get and print the messages
		var messages snmpsimclient.MessageMetrics
		messages, err = client.GetMessages(filters)
		if err != nil {
			log.Error().
				Msg("Error while getting of messages")
			os.Exit(1)
		}
		err = printData(messages, format, prettified, depth)
		if err != nil {
			log.Error().
				Msg("Error while printing data")
			os.Exit(1)
		}
	},
}

func init() {
	MessagesCmd.Flags().String("protocol", "", "Setting the protocol message filter")
	MessagesCmd.Flags().String("local_address", "", "Setting the local_address message filter")
	MessagesCmd.Flags().String("peer_address", "", "Setting the peer_address message filter")
	MessagesCmd.Flags().String("engine_id", "", "Setting the engine_id message filter")
	MessagesCmd.Flags().String("security_model", "", "Setting the security_model message filter")
	MessagesCmd.Flags().String("security_level", "", "Setting the security_level message filter")
	MessagesCmd.Flags().String("context_engine_id", "", "Setting the context_engine_id message filter")
	MessagesCmd.Flags().String("context_name", "", "Setting the context_name message filter")
	MessagesCmd.Flags().String("pdu_type", "", "Setting the pdu_type message filter")
	MessagesCmd.Flags().String("recording", "", "Setting the recording message filter")
}

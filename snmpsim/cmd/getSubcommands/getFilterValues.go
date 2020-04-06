package getsubcommands

import (
	snmpsimclient "github.com/inexio/snmpsim-restapi-go-client"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// GetFilterValuesCmd represents the getFilterValues command
var GetFilterValuesCmd = &cobra.Command{
	Use:   "filter-values",
	Args:  cobra.ExactArgs(0),
	Short: "Returns all possible values for a given filter",
	Long:  `Returns all possible filter values for either message or packet filters`,
	Run: func(cmd *cobra.Command, args []string) {
		//Parse all persistent flags
		format, depth, prettified := parsePersistentFlags(cmd)
		//check flags
		if !cmd.Flag("message").Changed && !cmd.Flag("packet").Changed {
			log.Error().
				Msg("Either 'message' or 'packet' flag has to be set")
			os.Exit(1)
		} else if cmd.Flag("message").Changed && cmd.Flag("packet").Changed {
			log.Error().
				Msg("Only 'message' or 'packet' flag can be set. Not both at the same time.")
			os.Exit(1)
		}

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

		//Read in the flag input
		if cmd.Flag("message").Changed {
			mFilter := cmd.Flag("message").Value.String()
			var mValues []string
			mValues, err = client.GetPossibleValuesForMessageFilter(mFilter)
			if err != nil {
				log.Error().
					Msg("Error while getting possible values for message filter")
				os.Exit(1)
			}

			err = printData(mValues, format, prettified, depth)
			if err != nil {
				log.Error().
					Msg("Error while printing data")
				os.Exit(1)
			}
		} else if cmd.Flag("packet").Changed {
			pFilter := cmd.Flag("packet").Value.String()
			var pValues []string
			pValues, err = client.GetPossibleValuesForMessageFilter(pFilter)
			if err != nil {
				log.Error().
					Msg("Error while getting possible values for packet filter")
				os.Exit(1)
			}

			err = printData(pValues, format, prettified, depth)
			if err != nil {
				log.Error().
					Msg("Error while printing data")
				os.Exit(1)
			}
		}
	},
}

func init() {
	GetFilterValuesCmd.Flags().String("message", "", "Returns all possible values for a given message filter")
	GetFilterValuesCmd.Flags().String("packet", "", "Returns all possible values for a given packet filter")
}

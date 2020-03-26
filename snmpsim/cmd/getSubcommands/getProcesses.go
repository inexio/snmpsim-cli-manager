package getsubcommands

import (
	"os"

	snmpsimclient "github.com/inexio/snmpsim-restapi-go-client"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ProcessesCmd represents the processes command
var ProcessesCmd = &cobra.Command{
	Use:   "processes",
	Args:  cobra.ExactArgs(0),
	Short: "Returns a list of all processes",
	Long: `Returns a complete list of all running processes.

All details of one specific process can be retrieved via 'get process <process-id>'.`,
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
				Msg("Error while creating new metrics client")
			os.Exit(1)
		}
		err = client.SetUsernameAndPassword(username, password)
		if err != nil {
			log.Error().
				Msg("Error while setting username of password")
			os.Exit(1)
		}

		//Get and print the processes
		var filters map[string]string
		var processes snmpsimclient.ProcessesMetrics
		processes, err = client.GetProcesses(filters)
		if err != nil {
			log.Error().
				Msg("Error while getting processes")
			os.Exit(1)
		}
		err = printData(processes, format, prettified, depth)
		if err != nil {
			log.Error().
				Msg("Error while printing data")
			os.Exit(1)
		}
	},
}

package getsubcommands

import (
	"os"
	"strconv"

	snmpsimclient "github.com/inexio/snmpsim-restapi-go-client"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ProcessCmd represents the process command
var ProcessCmd = &cobra.Command{
	Use:   "process <id>",
	Args:  cobra.ExactArgs(1),
	Short: "Returns details about a process",
	Long: `Returns details about the process with the given id.

These details include:
	- id
	- runtime
	- cpu
	- memory
	- files
	- exits
	- changes
	- last-update
	- console-pages
	- supervisor`,
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
				Msg("Error while creating new metrics client")
			os.Exit(1)
		}
		err = client.SetUsernameAndPassword(username, password)
		if err != nil {
			log.Error().
				Msg("Error while setting username and password")
			os.Exit(1)
		}

		//Read in and convert the process-id
		id, err := strconv.Atoi(args[0])
		if err != nil {
			log.Error().
				Msg("Error while converting id from string to integer")
			os.Exit(1)
		}

		//Initialize data var
		var data interface{}
		if cmd.Flag("console").Changed {
			var processConsolePages snmpsimclient.Consoles
			processConsolePages, err := client.GetProcessConsolePages(id)
			if err != nil {
				log.Error().
					Msg("Error while getting process console pages")
				os.Exit(1)
			}

			data = processConsolePages
		} else if cmd.Flag("console-page").Changed {
			//Read in page-id
			var pageID int
			pageID, err = cmd.Flags().GetInt("console-page")
			if err != nil {
				log.Error().
					Msg("Error while retrieving console-page id")
				os.Exit(1)
			}

			//Get process console page via page-id
			var processConsolePage snmpsimclient.Console
			processConsolePage, err = client.GetProcessConsolePage(id, pageID)
			if err != nil {
				log.Error().
					Msg("Error while getting process console page")
				os.Exit(1)
			}

			data = processConsolePage
		} else if cmd.Flag("endpoints").Changed {
			//Get process endpoints
			var processEndpoints snmpsimclient.ProcessEndpoints
			processEndpoints, err = client.GetProcessEndpoints(id)
			if err != nil {
				log.Error().
					Msg("Error while getting process endpoints")
				os.Exit(1)
			}

			data = processEndpoints
		} else if cmd.Flag("endpoint").Changed {
			//Read in endpoint-id
			endpointID, err := cmd.Flags().GetInt("endpoint")
			if err != nil {
				log.Error().
					Msg("Error while reading in endpoint-id")
				os.Exit(1)
			}

			//Get process endpoint via endpoint-id
			var processEndpoint snmpsimclient.ProcessEndpoint
			processEndpoint, err = client.GetProcessEndpoint(id, endpointID)
			if err != nil {
				log.Error().
					Msg("Error while getting process endpoint")
				os.Exit(1)
			}

			data = processEndpoint
		} else {
			//Get the process via id
			var process snmpsimclient.ProcessMetrics
			process, err = client.GetProcess(id)
			if err != nil {
				log.Error().
					Msg("Error while getting processes")
				os.Exit(1)
			}

			data = process
		}

		//Print the data
		err = printData(data, format, prettified, depth)
		if err != nil {
			log.Error().
				Msg("Error while printing data")
			os.Exit(1)
		}
	},
}

func init() {
	//Define console flags
	ProcessCmd.Flags().Int("console-page", 0, "Show the processes output of a certain console-page")
	ProcessCmd.Flags().Bool("console", false, "Show the processes console outputs")

	//Define endpoint flags
	ProcessCmd.Flags().Int("endpoint", 0, "Show details about the process-endpoint with the given id")
	ProcessCmd.Flags().Bool("endpoints", false, "Show all process-endpoints")
}

package getsubcommands

import (
	snmpsimclient "github.com/inexio/snmpsim-restapi-go-client"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// GetRecordFilesCmd represents the getRecordFiles command
var GetRecordFilesCmd = &cobra.Command{
	Use:   "record-files",
	Args:  cobra.ExactArgs(0),
	Short: "Returns all record files.",
	Long: `Returns a detailed list of all record files available.

All details of one specific record-file can be retrieved via 'get record-file <remote-path>'.`,
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
				Msg("Error while creating management client")
			os.Exit(1)
		}
		err = client.SetUsernameAndPassword(username, password)
		if err != nil {
			log.Error().
				Msg("Error while setting username and password")
			os.Exit(1)
		}

		//Get and print the record-files
		var recordFiles snmpsimclient.Recordings
		recordFiles, err = client.GetRecordFiles()
		if err != nil {
			log.Error().
				Msg("Error while getting record files")
			os.Exit(1)
		}
		err = printData(recordFiles, format, prettified, depth)
		if err != nil {
			log.Error().
				Msg("Error while printing data")
			os.Exit(1)
		}
	},
}

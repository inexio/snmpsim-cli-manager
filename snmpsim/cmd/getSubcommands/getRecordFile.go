package getsubcommands

import (
	snmpsimclient "github.com/inexio/snmpsim-restapi-go-client"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// GetRecordFileCmd represents the getRecordFile command
var GetRecordFileCmd = &cobra.Command{
	Use:   "record-file <remote path>",
	Args:  cobra.ExactArgs(1),
	Short: "Returns a record-file",
	Long:  `Returns the contents of the record-file that is located at the given path.`,
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

		//Read in the remote path
		path := args[0]

		//Get and print the record-file contents
		var recordFile string
		recordFile, err = client.GetRecordFile(path)
		if err != nil {
			log.Error().
				Msg("Error while getting record file")
			os.Exit(1)
		}
		err = printData(recordFile, format, prettified, depth)
		if err != nil {
			log.Error().
				Msg("Error while printing data")
			os.Exit(1)
		}
	},
}

package deletesubcommands

import (
	"fmt"
	snmpsimclient "github.com/inexio/snmpsim-restapi-go-client"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// DeleteRecordFileCmd represents the deleteRecordFile command
var DeleteRecordFileCmd = &cobra.Command{
	Use:   "record-file <path>",
	Args:  cobra.ExactArgs(1),
	Short: "Deletes a record-file",
	Long:  `Deletes the record file at the given path`,
	Run: func(cmd *cobra.Command, args []string) {
		//Load the client data from the config
		baseURL := viper.GetString("mgmt.http.baseURL")
		username := viper.GetString("mgmt.http.authUsername")
		password := viper.GetString("mgmt.http.authPassword")

		//Create a new client
		client, err := snmpsimclient.NewManagementClient(baseURL)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error while creating management client")
			os.Exit(1)
		}
		err = client.SetUsernameAndPassword(username, password)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error while setting username and password")
			os.Exit(1)
		}

		//Read in the record-path
		path := args[0]

		//Delete the record-file
		err = client.DeleteRecordFile(path)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error while deleting record file")
			os.Exit(1)
		}

		fmt.Println("Record file at", args[0], "has been deleted successfully.")
	},
}

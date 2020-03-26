package cmd

import (
	"fmt"
	"os"

	snmpsimclient "github.com/inexio/snmpsim-restapi-go-client"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// uploadRecordCmd represents the uploadRecord command
var uploadRecordCmd = &cobra.Command{
	Use:   "upload-record <remote path>",
	Args:  cobra.ExactArgs(1),
	Short: "Uploads a snmprec to the api.",
	Long: `Uploads a snmp record-file to the api-server.

You can either upload a file as a whole or use a string that will be written into a file.`,
	Run: func(cmd *cobra.Command, args []string) {
		//Check flags
		if !cmd.Flag("file").Changed && !cmd.Flag("string").Changed {
			log.Error().
				Msg("Either --file or --string has to be set")
			os.Exit(1)
		} else if cmd.Flag("file").Changed && cmd.Flag("string").Changed {
			log.Error().
				Msg("Only --file or --string can be set. Not both at the same time.")
			os.Exit(1)
		}

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

		//Read in flag values
		rPath := args[0]
		lPath := cmd.Flag("file").Value.String()
		fileString := cmd.Flag("string").Value.String()

		//Upload a file or a string
		if lPath != "" {
			err = client.UploadRecordFile(lPath, rPath)
			if err != nil {
				log.Error().
					Msg("Error while uploading the record file")
				os.Exit(1)
			}
			fmt.Println("Record file successfully uploaded from", lPath, "to", rPath)
		} else if fileString != "" {
			err = client.UploadRecordFileString(&fileString, rPath)
			if err != nil {
				log.Error().
					Msg("Error while uploading the string into the record file")
				os.Exit(1)
			}
			fmt.Println("String has been successfully uploaded into", rPath)
		}
	},
}

func init() {
	rootCmd.AddCommand(uploadRecordCmd)
	uploadRecordCmd.Flags().String("file", "", "Path of file which is to be uploaded")
	uploadRecordCmd.Flags().String("string", "", "Content of the file which is to be uploaded")
}

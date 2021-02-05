package getsubcommands

import (
	snmpsimclient "github.com/inexio/snmpsim-restapi-go-client"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strconv"
)

// GetEngineCmd represents the getEngine command
var GetEngineCmd = &cobra.Command{
	Use:   "engine <id>",
	Args:  cobra.ExactArgs(1),
	Short: "Returns details about an engine",
	Long: `Returns details about the engine with the given id.

These details include:
	- id
	- name
	- engine-id
	- endpoints
	- users
	- tags`,
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

		//Read in and convert the id
		id, err := strconv.Atoi(args[0])
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error during conversion of id from string to integer")
			os.Exit(1)
		}

		//Get and print the engine
		var engine snmpsimclient.Engine
		engine, err = client.GetEngine(id)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error while getting engine")
			os.Exit(1)
		}
		err = printData(engine, format, prettified, depth)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error while printing data")
			os.Exit(1)
		}
	},
}

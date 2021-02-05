package getsubcommands

import (
	snmpsimclient "github.com/inexio/snmpsim-restapi-go-client"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strconv"
)

// GetUserCmd represents the getUser command
var GetUserCmd = &cobra.Command{
	Use:   "user <id>",
	Args:  cobra.ExactArgs(1),
	Short: "Returns details about a user",
	Long: `Returns details about the user with the given id.

These details include:
	- id
	- name
	- auth-key
	- auth-proto
	- priv-key
	- priv-proto
	- username
	- tag`,
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

		//Get and print the user
		var user snmpsimclient.User
		user, err = client.GetUser(id)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error while getting user")
			os.Exit(1)
		}
		err = printData(user, format, prettified, depth)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error while printing data")
			os.Exit(1)
		}
	},
}

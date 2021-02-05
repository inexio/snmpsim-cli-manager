package getsubcommands

import (
	snmpsimclient "github.com/inexio/snmpsim-restapi-go-client"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// GetUsersCmd represents the getUsers command
var GetUsersCmd = &cobra.Command{
	Use:   "users",
	Args:  cobra.ExactArgs(0),
	Short: "Returns a list of all users",
	Long: `Returns a detailed list of all users available.

All details of one specific user can be retrieved via 'get user <user-id>'.`,
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
		if username != "" && password != "" {
			err = client.SetUsernameAndPassword(username, password)
			if err != nil {
				log.Error().
					Err(err).
					Msg("Error while setting username and password")
				os.Exit(1)
			}
		}

		//Parse filters from flags
		filters := parseFilters(cmd)

		//Get and print the users
		var users snmpsimclient.Users
		users, err = client.GetUsers(filters)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error while getting users")
			os.Exit(1)
		}
		err = printData(users, format, prettified, depth)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error while printing data")
			os.Exit(1)
		}
	},
}

func init() {
	GetUsersCmd.Flags().String("name", "", "Set the name filter")
	GetUsersCmd.Flags().String("user", "", "Set the user filter")
	GetUsersCmd.Flags().String("auth_proto", "", "Set the authentication protocol filter")
	GetUsersCmd.Flags().String("priv_proto", "", "Set the private protocol filter")
}

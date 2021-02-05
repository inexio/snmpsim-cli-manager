package addsubcommands

import (
	"fmt"
	"os"

	snmpsimclient "github.com/inexio/snmpsim-restapi-go-client"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// UserToEngineCmd represents the userToEngine command
var UserToEngineCmd = &cobra.Command{
	Use:   "user-to-engine",
	Args:  cobra.ExactArgs(0),
	Short: "Adds an user to an engine",
	Long:  `Adds the user with a given user-id to the engine with the given engine-id.`,
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

		//Read in the engine-id
		engineID, err := cmd.Flags().GetInt("engine")
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error while retrieving engine-id")
			os.Exit(1)
		}

		//Read in the user-id
		userID, err := cmd.Flags().GetInt("user")
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error while retrieving user-id")
			os.Exit(1)
		}

		//Add the user to the engine
		err = client.AddUserToEngine(engineID, userID)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error while adding user to engine")
			os.Exit(1)
		}
		fmt.Println("User", userID, "has been added to engine", engineID)
	},
}

func init() {
	//Set user flag
	UserToEngineCmd.Flags().Int("user", 0, "Id of the user that is to be added to the engine")
	err := UserToEngineCmd.MarkFlagRequired("user")
	if err != nil {
		log.Error().
			Err(err).
			Msg("Could not mark 'user' flag required")
		os.Exit(1)
	}

	//Set engine flag
	UserToEngineCmd.Flags().Int("engine", 0, "Id of the engine to that the user will be added")
	err = UserToEngineCmd.MarkFlagRequired("engine")
	if err != nil {
		log.Error().
			Err(err).
			Msg("Could not mark 'engine' flag required")
		os.Exit(1)
	}
}

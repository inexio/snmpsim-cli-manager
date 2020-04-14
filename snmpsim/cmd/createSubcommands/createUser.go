package createsubcommands

import (
	"fmt"
	snmpsimclient "github.com/inexio/snmpsim-restapi-go-client"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// CreateUserCmd represents the createUser command
var CreateUserCmd = &cobra.Command{
	Use:   "user",
	Args:  cobra.ExactArgs(0),
	Short: "Creates a new user",
	Long:  `Creates a new user and returns its id.`,
	Run: func(cmd *cobra.Command, args []string) {
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

		//Read in the users name, authKey, authProto, privKey, privProto and userFlag
		userFlag := cmd.Flag("user").Value.String()
		name := cmd.Flag("name").Value.String()
		authKey := cmd.Flag("authKey").Value.String()
		authProto := cmd.Flag("authProto").Value.String()
		privKey := cmd.Flag("privKey").Value.String()
		privProto := cmd.Flag("privProto").Value.String()

		//creating the user
		var user snmpsimclient.User
		if cmd.Flag("tag").Changed {
			//Read in tag-id
			tagID, err := cmd.Flags().GetInt("tag")
			if err != nil {
				log.Error().
					Msg("Error while retrieving tagID")
				os.Exit(1)
			}

			//Validate tag-id
			if tagID == 0 {
				log.Error().
					Msg("tagID can not be 0")
				os.Exit(1)
			}

			//Check if tag with given id exists
			_, err = client.GetTag(tagID)
			if err != nil {
				log.Error().
					Msg("No tag with the given id found")
				os.Exit(1)
			}

			user, err = client.CreateUserWithTag(userFlag, name, authKey, authProto, privKey, privProto, tagID)
			if err != nil {
				log.Error().
					Msg("Error during creation of the user")
				os.Exit(1)
			}
		} else {
			user, err = client.CreateUser(userFlag, name, authKey, authProto, privKey, privProto)
			if err != nil {
				log.Error().
					Msg("Error during creation of the user")
				os.Exit(1)
			}
		}

		fmt.Println("User has been created successfully.")
		fmt.Println("Id:", user.ID)

		//Add user to engine (if engine flag is set)
		if cmd.Flag("engine").Changed {
			//Read in engine-id
			engineID, err := cmd.Flags().GetInt("engine")
			if err != nil {
				log.Error().
					Msg("Error while retrieving engineID")
				os.Exit(1)
			}

			//Check if engine with given id exists
			_, err = client.GetEngine(engineID)
			if err != nil {
				log.Error().
					Msg("No engine with the given id found")
				os.Exit(1)
			}

			//Add user to engine
			err = client.AddUserToEngine(engineID, user.ID)
			if err != nil {
				log.Error().
					Msg("Error while adding user to engine")
				os.Exit(1)
			}
			fmt.Println("Successfully added user", user.ID, "to engine", engineID)
		}
	},
}

func init() {
	CreateUserCmd.Flags().String("user", "", "The user of the user")
	err := CreateUserCmd.MarkFlagRequired("user")
	if err != nil {
		log.Error().
			Msg("Could not mark 'user' flag required")
		os.Exit(1)
	}

	CreateUserCmd.Flags().String("authKey", "", "The authentication key of the user")
	err = CreateUserCmd.MarkFlagRequired("authKey")
	if err != nil {
		log.Error().
			Msg("Could not mark 'authKey' flag required")
		os.Exit(1)
	}

	CreateUserCmd.Flags().String("authProto", "", "The authentication protocol of the user")
	err = CreateUserCmd.MarkFlagRequired("authProto")
	if err != nil {
		log.Error().
			Msg("Could not mark 'authProto' flag required")
		os.Exit(1)
	}

	CreateUserCmd.Flags().String("privKey", "", "The private key of the user")
	err = CreateUserCmd.MarkFlagRequired("privKey")
	if err != nil {
		log.Error().
			Msg("Could not mark 'privKey' flag required")
		os.Exit(1)
	}

	CreateUserCmd.Flags().String("privProto", "", "The private protocol of the user")
	err = CreateUserCmd.MarkFlagRequired("privProto")
	if err != nil {
		log.Error().
			Msg("Could not mark 'privProto' flag required")
		os.Exit(1)
	}

	CreateUserCmd.Flags().Int("tag", 0, "Optional flag to mark an user with a tag (requires a tag-id)")
	CreateUserCmd.Flags().Int("engine", 0, "Optional flag to add the user to an engine (requires an engine-id)")
}

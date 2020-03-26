package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	snmpsimclient "github.com/inexio/snmpsim-restapi-go-client"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// eraseEnvCmd represents the eraseEnv command
var eraseEnvCmd = &cobra.Command{
	Use:   "erase-env <tag-id>",
	Args:  cobra.ExactArgs(1),
	Short: "Erases an lab environment",
	Long:  `Completely deletes all components created during setup-env operation including the tag created with it.`,
	Run: func(cmd *cobra.Command, args []string) {
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

		//Create reader
		reader := bufio.NewReader(os.Stdin)

		//Read in the tag-id
		tagId, err := strconv.Atoi(args[0])
		if err != nil {
			log.Error().
				Msg("Error while converting " + args[0] + "from string to int")
			os.Exit(1)
		}

		//Get information about the tag
		tag, err := client.GetTag(tagId)
		if err != nil {
			log.Error().
				Msg("Error while getting tag")
			os.Exit(1)
		}

		deleteEnv := true

		//Check if the force flag is set
		if !cmd.Flag("force").Changed {
			fmt.Print("Are you sure you to delete the environment tagged with ", tag.Name, " id: ", tagId, "?(yes/no) ")
			//checking the user inpit
			input, err := reader.ReadString('\n')
			if err != nil {
				log.Error().
					Msg("Error while retrieving input")
				os.Exit(1)
			}

			//Remove carriage return and line feed characters
			input = strings.Replace(strings.Replace(input, "\n", "", -1), "\r", "", -1)

			switch input {
			case "yes", "Yes", "y", "Y":
				deleteEnv = true
			case "no", "No", "n", "N":
				deleteEnv = false
			default:
				log.Debug().
					Msg("Invalid input: " + input)
				os.Exit(1)
			}
		}

		if deleteEnv {
			//Delete all tagged objects
			_, err = client.DeleteAllObjectsWithTag(tagId)
			if err != nil {
				log.Error().
					Msg("Error while deleting all objects tagged with " + tag.Name)
				os.Exit(1)
			}

			//Delete the tag itself
			err = client.DeleteTag(tagId)
			if err != nil {
				log.Error().
					Msg("Error while deleting tag " + tag.Name)
				os.Exit(1)
			}
			fmt.Println("Environment", tag.Name, "id", tagId, "has been deleted successfully.")
		}
	},
}

func init() {
	rootCmd.AddCommand(eraseEnvCmd)
	eraseEnvCmd.Flags().BoolP("force", "f", false, "Disables the 'Are you sure you want to delete this' question")
}

package getsubcommands

import (
	//"github.com/davecgh/go-spew/spew"
	snmpsimclient "github.com/inexio/snmpsim-restapi-go-client"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// GetLabsCmd represents the getLabs command
var GetLabsCmd = &cobra.Command{
	Use:   "labs",
	Args:  cobra.ExactArgs(0),
	Short: "Returns a list of all labs",
	Long: `Returns a detailed list of all labs available.

All details of one specific lab can be retrieved via 'get lab <lab-id>'.`,
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

		//Parse filters from flags
		filters := parseFilters(cmd)

		//Get and print the labs
		var labs snmpsimclient.Labs
		labs, err = client.GetLabs(filters)
		if err != nil {
			log.Error().
				Msg("Error while getting labs")
			os.Exit(1)
		}
		err = printData(labs, format, prettified, depth)
		if err != nil {
			log.Error().
				Msg("Error while printing data")
			os.Exit(1)
		}
	},
}

func init() {
	GetLabsCmd.Flags().String("name", "", "Set the name filter")
	GetLabsCmd.Flags().String("power", "", "Set the power filter")
}

package getsubcommands

import (
	snmpsimclient "github.com/inexio/snmpsim-restapi-go-client"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// GetTagsCmd represents the getTags command
var GetTagsCmd = &cobra.Command{
	Use:   "tags",
	Args:  cobra.ExactArgs(0),
	Short: "Returns a list of all tags.",
	Long: `Returns a detailed list of all tags available.

All details of one specific tag can be retrieved via 'get tag <tag-id>'.`,
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

		//Get and print the tags
		var tags snmpsimclient.Tags
		tags, err = client.GetTags(filters)
		if err != nil {
			log.Error().
				Msg("Error while getting tags")
			os.Exit(1)
		}
		err = printData(tags, format, prettified, depth)
		if err != nil {
			log.Error().
				Msg("Error while printing data")
			os.Exit(1)
		}
	},
}

func init() {
	GetTagsCmd.Flags().String("name", "", "Set the name filter")
}

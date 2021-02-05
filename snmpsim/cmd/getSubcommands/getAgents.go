package getsubcommands

import (
	snmpsimclient "github.com/inexio/snmpsim-restapi-go-client"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// GetAgentsCmd represents the getAgents command
var GetAgentsCmd = &cobra.Command{
	Use:   "agents",
	Args:  cobra.ExactArgs(0),
	Short: "Returns a list of all agents",
	Long: `Returns a detailed list of all agents available.

All details of one specific agent can be retrieved via 'get agent <agent-id>'.`,
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

		//Parse filters from flags
		filters := parseFilters(cmd)

		//Get and print the agents
		var agents snmpsimclient.Agents
		agents, err = client.GetAgents(filters)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error while getting agents")
			os.Exit(1)
		}
		err = printData(agents, format, prettified, depth)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error while printing data")
			os.Exit(1)
		}
	},
}

func init() {
	GetAgentsCmd.Flags().String("name", "", "Set the name filter")
	GetAgentsCmd.Flags().String("data_dir", "", "Set the data directory filter")
}

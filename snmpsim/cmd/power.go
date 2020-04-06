package cmd

import (
	"fmt"
	"os"
	"strconv"

	snmpsimclient "github.com/inexio/snmpsim-restapi-go-client"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

// powerCmd represents the power command
var powerCmd = &cobra.Command{
	Use:   "power <lab_id> --on/--off",
	Args:  cobra.ExactArgs(1),
	Short: "Sets the power of a lab",
	Long:  `Sets the power setting of a lab to either 'on' or 'off'"`,
	Run: func(cmd *cobra.Command, args []string) {
		//Check flags
		if !cmd.Flag("on").Changed && !cmd.Flag("off").Changed {
			log.Error().
				Msg("Either --on or --off has to be set")
			os.Exit(1)
		} else if cmd.Flag("on").Changed && cmd.Flag("off").Changed {
			log.Error().
				Msg("Only --on or --of can be set. Not both at the same time.")
			os.Exit(1)
		}

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

		//Read in the lab-id
		labID := args[0]
		lab, err := strconv.Atoi(labID)
		if err != nil {
			log.Error().
				Msg("Error while converting labID from string to integer")
			os.Exit(1)
		}

		//Read in flag vlaues
		on, err := cmd.Flags().GetBool("on")
		if err != nil {
			log.Error().
				Msg("Error while retrieving --on flag value")
			os.Exit(1)
		}
		off, err := cmd.Flags().GetBool("off")
		if err != nil {
			log.Error().
				Msg("Error while retrieving --off flag value")
			os.Exit(1)
		}

		//Convert the flag value
		var power bool
		var status string
		if on {
			power = true
			status = "on"
		} else if off {
			power = false
			status = "off"
		}

		//Set the labs power
		err = client.SetLabPower(lab, power)
		if err != nil {
			log.Error().
				Msg("Error during setting of lab power")
			os.Exit(1)
		}
		fmt.Println("Power of Lab", labID, "has been successfully set to", status)
	},
}

func init() {
	rootCmd.AddCommand(powerCmd)
	powerCmd.Flags().Bool("on", false, "Set the labs power to 'on'")
	powerCmd.Flags().Bool("off", false, "Set the labs power to 'off'")
}

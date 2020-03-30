package setsubcommands

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// SetEnvConfigCmd represents the setEnvConfig command
var SetEnvConfigCmd = &cobra.Command{
	Use:   "env-config",
	Args:  cobra.ExactArgs(1),
	Short: "Sets the path of the env-config file",
	Long:  `Permanently sets the path of the env-config file`,
	Run: func(cmd *cobra.Command, args []string) {
		cfgFile := args[0]

		cfgManager := viper.New()
		cfgManager.SetConfigFile("config/cfg-mgmt.yaml")
		if err := cfgManager.ReadInConfig(); err != nil {
			log.Debug().
				Msg("Could not read in cfg-mgmt file")
		}
		cfgManager.Set("env-config", cfgFile)
		cfgManager.WriteConfig()
	},
}

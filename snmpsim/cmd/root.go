package cmd

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
	"strings"

	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "snmpsim",
	Short: "A cli client for managing lab environments created by snmpsim",
	Long:  `This cli management tool gives you the opportunity to fully work with snmpsim via the command line.`,
}

/*
Execute represents the entry point of the application
*/
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	//set log output format
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/go/src/snmpsim-cli-manager/snmpsim/config/snmpsim-cli-manager-config.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	//Set env var prefix to only match certain vars
	viper.SetEnvPrefix("SNMPSIM_CLI")
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()

	if cfgFile == "" {
		cfgFile = viper.GetString("config")
	}

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("config/")
		viper.SetConfigType("yaml")
		viper.SetConfigName("snmpsim-cli-manager-config")
	}

	err := viper.ReadInConfig()
	if err != nil {
		log.Debug().
			Msg("Could not read in config file\nMake sure you have set environment variables if you're not using a config-file")
	}
}

package cmd

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/speckleworks/gospeckle/pkg"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var contextName string
var currentConfig CurrentConfig
var speckleClient gospeckle.Client
var ctx context.Context

// var defaultConfig = map[string]string{"tag": "tags", "category": "categories"}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gospeckle/config.yaml)")
	rootCmd.PersistentFlags().StringVarP(&contextName, "context", "c", "", "configuration server/user context")

	// viper.BindPFlag("context", rootCmd.PersistentFlags().Lookup("context"))
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var defaultCfgDir = string(home + "/.gospeckle")
	var defaultCfgPath = string(defaultCfgDir + "/config.yaml")

	// Don't forget to read config either from cfgFile or from home directory!
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {

		cfgFile = defaultCfgPath

		// Search config in home directory with name ".gospeckle".
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(home + "/.gospeckle")
		viper.AddConfigPath(".")
	}

	if err := viper.ReadInConfig(); err != nil {
		if cfgFile == defaultCfgPath {
			err := writeDefaultConfig(defaultCfgDir, cfgFile)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			fmt.Println("Can't read config:", err)
			os.Exit(1)
		}
	}
}

var rootCmd = &cobra.Command{
	Use:   "gospeckle",
	Short: "GoSpeckle is a command line client to interact with speckle servers",
	Long:  `GoSpeckle is a command line client to interact with speckle servers`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		ctx = context.TODO()
		currentConfig = getConfigContext()
		httpClient := new(http.Client)
		host, err := url.Parse(currentConfig.Server.Host)

		if err != nil {
			fmt.Println("Could not parse current context server host url")
			fmt.Println(err)
			os.Exit(1)
		}

		speckleClient = *gospeckle.NewClient(httpClient, host, nil, currentConfig.Server.Version, currentConfig.User.Token)

	},
	Run: func(cmd *cobra.Command, args []string) {
		printLogo()
		cmd.Help()
	},
}

// Execute is the entrypoint for the cmd package
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

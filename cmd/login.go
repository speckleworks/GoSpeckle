package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var password string

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringVarP(&password, "password", "p", "", "The password for the user in the specified or current config")
	loginCmd.MarkFlagRequired("password")
}

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Save an authentication token to your current context config.",
	Run: func(cmd *cobra.Command, args []string) {
		err := speckleClient.Login(ctx, currentConfig.User.Email, password, true)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		config, err := getConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		user, err := config.getUserByName(currentConfig.User.Name)

		user.Token = speckleClient.Token

		viper.Set("users", config.Users)
		viper.WriteConfig()
		fmt.Printf("Succesfully logged into %v with user email %v", currentConfig.Server.Host, currentConfig.User.Email)
	},
}

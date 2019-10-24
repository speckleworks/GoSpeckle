package cmd

import (
	"fmt"
	"os"

	gospeckle "github.com/speckleworks/gospeckle/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var userName string
var userSurname string
var company string
var password string

func init() {
	rootCmd.AddCommand(accountCmd)

	accountCmd.AddCommand(loginCmd)
	accountCmd.AddCommand(registerCmd)

	loginCmd.Flags().StringVarP(&password, "password", "p", "", "The password for the user in the specified or current config")
	loginCmd.MarkFlagRequired("password")

	registerCmd.Flags().StringVar(&userName, "name", "", "name of the user")
	registerCmd.Flags().StringVar(&userSurname, "surname", "", "surname of the user")
	registerCmd.Flags().StringVar(&email, "email", "", "user's email")
	registerCmd.Flags().StringVar(&company, "company", "", "company the user works for")
	registerCmd.Flags().StringVar(&password, "password", "", "new user's password")
	registerCmd.MarkFlagRequired("email")
	registerCmd.MarkFlagRequired("password")
}

// loginCmd represents the login command
var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "commands to manage accounts specifically",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
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

// loginCmd represents the login command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a new user to the current config's server and set as new context",
	Run: func(cmd *cobra.Command, args []string) {
		newUser := gospeckle.AccountRegisterRequest{
			Name:     userName,
			Surname:  userSurname,
			Email:    email,
			Company:  company,
			Password: password,
		}

		err := speckleClient.Account.Register(ctx, newUser)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = speckleClient.Login(ctx, email, password, true)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		config, err := getConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		newConfigUser := ConfigUser{
			Name:  userName,
			Email: email,
			Token: speckleClient.Token,
		}

		serverName := currentConfig.Server.Name

		newConfigContext := ConfigContext{
			Name:   serverName + "-" + newConfigUser.Name,
			User:   newConfigUser.Name,
			Server: serverName,
		}

		config.Users = append(config.Users, newConfigUser)
		config.Contexts = append(config.Contexts, newConfigContext)

		viper.Set("users", config.Users)
		viper.Set("contexts", config.Contexts)
		viper.Set("current-context", newConfigContext.Name)

		viper.WriteConfig()
		fmt.Printf("Succesfully registered %v to server %v", newConfigUser.Name, currentConfig.Server.Host)
	},
}

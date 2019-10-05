package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	user          string
	email         string
	token         string
	server        string
	host          string
	serverVersion string
)

func init() {
	getConfigUsersCmd.Flags().StringVarP(&user, "name", "n", "", "the name of the user to be retrieved")
	getConfigServersCmd.Flags().StringVarP(&server, "name", "n", "", "the name of the server to be retrieved")
	getConfigContextsCmd.Flags().StringVarP(&contextName, "name", "n", "", "the name of the context to be retrieved")

	setConfigUserCmd.Flags().StringVarP(&user, "name", "n", "", "the name of the user to be updated/created")
	setConfigUserCmd.MarkFlagRequired("name")
	setConfigUserCmd.Flags().StringVarP(&email, "email", "e", "", "the email of the user to be updated/created")
	setConfigUserCmd.Flags().StringVarP(&token, "token", "t", "", "the API access token of the user to be updated/created")

	setConfigServerCmd.Flags().StringVarP(&server, "name", "n", "", "the name of the server to be updated/created")
	setConfigServerCmd.MarkFlagRequired("name")
	setConfigServerCmd.Flags().StringVar(&host, "host", "", "the http url of the server to be updated/created")
	setConfigServerCmd.Flags().StringVar(&serverVersion, "version", "", "the api version of the server to be updated/created")

	setConfigContextCmd.Flags().StringVarP(&contextName, "name", "n", "", "the name of the context to be updated/created")
	setConfigContextCmd.MarkFlagRequired("name")
	setConfigContextCmd.Flags().StringVar(&server, "server", "", "the new context server name")
	setConfigContextCmd.Flags().StringVar(&user, "user", "", "the new context user name")

	useConfigContextCmd.Flags().StringVarP(&contextName, "name", "n", "", "the name of the context to be used")
	useConfigContextCmd.MarkFlagRequired("name")

	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(getConfigCmd)
	configCmd.AddCommand(currentContextCmd)
	configCmd.AddCommand(getConfigUsersCmd)
	configCmd.AddCommand(getConfigServersCmd)
	configCmd.AddCommand(getConfigContextsCmd)
	configCmd.AddCommand(setConfigUserCmd)
	configCmd.AddCommand(setConfigServerCmd)
	configCmd.AddCommand(setConfigContextCmd)
	configCmd.AddCommand(useConfigContextCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration management commands to manipulate the .gospeckle config file holding.",
	Long: `Configuration management commands to manipulate the .gospeckle config file holding.
	This file contains login information for different speckle servers.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
	},
}

var getConfigCmd = &cobra.Command{
	Use:   "get",
	Short: "Display the current config",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := getConfig()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = printYaml(c)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var currentContextCmd = &cobra.Command{
	Use:   "current-context",
	Short: "Display the current context",
	Run: func(cmd *cobra.Command, args []string) {
		c := getConfigContext()

		err := printYaml(c)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var getConfigUsersCmd = &cobra.Command{
	Use:   "get-user",
	Short: "Display the current config users",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := getConfig()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if user != "" {
			object, err := c.getUserByName(user)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			err = printYaml(object)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			err = printYaml(c.Users)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	},
}

var getConfigServersCmd = &cobra.Command{
	Use:   "get-server",
	Short: "Display the current config servers",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := getConfig()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if server != "" {
			object, err := c.getServerByName(server)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			err = printYaml(object)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			err = printYaml(c.Servers)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	},
}

var getConfigContextsCmd = &cobra.Command{
	Use:   "get-context",
	Short: "Display the current config contexts",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := getConfig()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if contextName != "" {
			object, err := c.getContextByName(contextName)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			err = printYaml(object)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			err = printYaml(c.Contexts)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	},
}

var setConfigUserCmd = &cobra.Command{
	Use:   "set-user",
	Short: "Update or create a new user in the current config file",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := getConfig()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		object, err := c.getUserByName(user)

		if object != nil {
			if email != "" {
				object.Email = email
			}
			if token != "" {
				object.Token = token
			}
			printYaml(object)

		} else {
			if email == "" {
				fmt.Println("Cannot create new user in config with empty email.")
				os.Exit(1)
			}
			newUser := ConfigUser{
				Name:  user,
				Email: email,
				Token: token,
			}
			printYaml(newUser)

			c.Users = append(c.Users, newUser)
		}

		viper.Set("users", c.Users)
		viper.WriteConfig()

	},
}

var setConfigServerCmd = &cobra.Command{
	Use:   "set-server",
	Short: "Update or create a new server in the current config file",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := getConfig()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		object, err := c.getServerByName(server)

		if object != nil {
			if host != "" {
				object.Host = host
			}
			if serverVersion != "" {
				object.Version = serverVersion
			}
			printYaml(object)
		} else {
			if host == "" {
				fmt.Println("Cannot create new server in config with empty host path.")
				os.Exit(1)
			}
			newServer := ConfigServer{
				Name: server,
				Host: host,
			}
			printYaml(newServer)
			c.Servers = append(c.Servers, newServer)
		}

		viper.Set("servers", c.Servers)
		viper.WriteConfig()

	},
}

var setConfigContextCmd = &cobra.Command{
	Use:   "set-context",
	Short: "Update or create a new context in the current config file",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := getConfig()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if user != "" {
			userObject, _ := c.getUserByName(user)

			if userObject == nil {
				fmt.Println("user does not exist in config so cannot be assigned to context")
				os.Exit(1)
			}
		}

		if server != "" {
			serverObject, _ := c.getServerByName(server)

			if serverObject == nil {
				fmt.Println("server does not exist in config so cannot be assigned to context")
				os.Exit(1)
			}
		}

		object, err := c.getContextByName(contextName)

		if object != nil {
			if user != "" {
				object.User = user
			}

			if server != "" {
				object.Server = server
			}
			printYaml(object)
		} else {

			if server == "" || user == "" {
				fmt.Println("Cannot create new context in config with empty server or user.")
				os.Exit(1)
			}

			newContext := ConfigContext{
				Name:   contextName,
				Server: server,
				User:   user,
			}
			printYaml(newContext)
			c.Contexts = append(c.Contexts, newContext)
		}

		viper.Set("contexts", c.Contexts)
		viper.WriteConfig()

	},
}

var useConfigContextCmd = &cobra.Command{
	Use:   "use-context",
	Short: "Set the current-context in the gospeckle config file",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := getConfig()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		object, err := c.getContextByName(contextName)

		if object == nil {
			fmt.Println("context does not exist in config so cannot be assigned to current-context")
			os.Exit(1)
		}

		c.CurrentContext = contextName

		viper.Set("current-context", contextName)
		viper.WriteConfig()

		printYaml(c.CurrentContext)
	},
}

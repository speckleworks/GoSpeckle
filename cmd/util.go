package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"
)

type ConfigContext struct {
	Name   string `yaml:"name"`
	User   string `yaml:"user"`
	Server string `yaml:"server"`
}

type ConfigUser struct {
	Name  string `yaml:"name"`
	Email string `yaml:"email"`
	Token string `yaml:"token"`
}

type ConfigServer struct {
	Name    string `yaml:"name"`
	Host    string `yaml:"host"`
	Version string `yaml:"version"`
}

type CurrentConfig struct {
	Name   string       `yaml:"name"`
	User   ConfigUser   `yaml:"user"`
	Server ConfigServer `yaml:"server"`
}

type Config struct {
	CurrentContext string          `yaml:"current-context" mapstructure:"current-context"`
	Servers        []ConfigServer  `yaml:"servers"`
	Users          []ConfigUser    `yaml:"users"`
	Contexts       []ConfigContext `yaml:"contexts"`
}

func printLogo() {
	logo := `
  _______     ______         ________  _______    _______   ______   __   ___  ___       _______  
 /" _   "|   /    " \       /"       )|   __ "\  /"     "| /" _  "\ |/"| /  ")|"  |     /"     "| 
(: ( \___)  // ____  \     (:   \___/ (. |__) :)(: ______)(: ( \___)(: |/   / ||  |    (: ______) 
 \/ \      /  /    ) :)     \___  \   |:  ____/  \/    |   \/ \     |    __/  |:  |     \/    |   
 //  \ ___(: (____/ //       __/  \\  (|  /      // ___)_  //  \ _  (// _  \   \  |___  // ___)_  
(:   _(  _|\        /       /" \   :)/|__/ \    (:      "|(:   _) \ |: | \  \ ( \_|:  \(:      "| 
 \_______)  \"_____/       (_______/(_______)    \_______) \_______)(__|  \__) \_______)\_______) 
																																																	
 `
	fmt.Print(logo)
}

func printYaml(y interface{}) error {
	p, err := yaml.Marshal(y)

	if err != nil {
		return err
	}

	fmt.Print(string(p))
	return nil
}

func printJSON(j interface{}) error {
	var prettyJSON bytes.Buffer
	jsonData, _ := json.Marshal(j)

	error := json.Indent(&prettyJSON, jsonData, "", "  ")
	if error != nil {
		return error
	}

	fmt.Println(string(prettyJSON.Bytes()))

	return nil
}

func (c Config) getUserByName(name string) (*ConfigUser, error) {
	sliceIndex := -1

	for i, u := range c.Users {
		if u.Name == name {
			sliceIndex = i
		}
	}

	if sliceIndex == -1 {
		return nil, fmt.Errorf("Could not find object with name: %v", name)
	}

	return &c.Users[sliceIndex], nil
}

func (c Config) getServerByName(name string) (*ConfigServer, error) {
	sliceIndex := -1

	for i, u := range c.Servers {
		if u.Name == name {
			sliceIndex = i
		}
	}

	if sliceIndex == -1 {
		return nil, fmt.Errorf("Could not find object with name: %v", name)
	}

	return &c.Servers[sliceIndex], nil
}

func (c Config) getContextByName(name string) (*ConfigContext, error) {
	sliceIndex := -1

	for i, u := range c.Contexts {
		if u.Name == name {
			sliceIndex = i
		}
	}

	if sliceIndex == -1 {
		return nil, fmt.Errorf("Could not find object with name: %v", name)
	}

	return &c.Contexts[sliceIndex], nil
}

var defaultConfig = []byte(`
current-context: default
servers:
	- name: hestia
		host: https://hestia.speckle.works
    version: v1
users:
  - name: anonymous
    email: ""
    token: ""
contexts:
  - name: default
    server: hestia
    user: anonymous
`)

func writeDefaultConfig(cfgDir string, cfgFile string) error {
	viper.ReadConfig(bytes.NewBuffer(defaultConfig))
	_, err := os.Stat(cfgDir)

	if os.IsNotExist(err) {
		err = os.MkdirAll(cfgDir, 0755)
		if err != nil {
			return err
		}

	}
	err = viper.WriteConfigAs(cfgFile)
	if err != nil {
		return err
	}
	fmt.Println("Wrote default config to " + cfgFile)
	return nil
}

func getConfig() (Config, error) {
	var c Config

	err := viper.Unmarshal(&c)

	return c, err
}

func getConfigContext() CurrentConfig {
	var config CurrentConfig
	var contexts []ConfigContext
	var users []ConfigUser
	var servers []ConfigServer

	if contextName == "" {
		contextName = viper.GetString("current-context")
	}

	config.Name = contextName

	viper.UnmarshalKey("contexts", &contexts)
	viper.UnmarshalKey("users", &users)
	viper.UnmarshalKey("servers", &servers)

	for _, c := range contexts {
		if c.Name == contextName {
			for _, s := range servers {
				if s.Name == c.Server {
					config.Server = s
				}
			}

			for _, u := range users {
				if u.Name == c.User {
					config.User = u
				}
			}
		}
	}

	return config
}

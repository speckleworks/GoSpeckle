package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/mitchellh/mapstructure"
	gospeckle "github.com/speckleworks/gospeckle/pkg"
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

type FileInput struct {
	Type     string                    `json:"type" yaml:"type"`
	Metadata gospeckle.RequestMetadata `json:"metadata" yaml:"metadata"`
	Spec     interface{}               `json:"spec" yaml:"spec"`
}

type Decoder interface {
	Decode(interface{}) error
}

type RequestObjects struct {
	Projects []gospeckle.ProjectRequest
	Streams  []gospeckle.StreamRequest
	Clients  []gospeckle.APIClientRequest
}

func (r RequestObjects) MakeRequests(ctx context.Context, c *gospeckle.Client) {
	var totalObjects []interface{}

	totalObjects = append(totalObjects, r.Projects, r.Streams, r.Clients)
	//  len(r.Projects) + len(r.Streams) + len(r.Clients)
	ch := make(chan string)

	for _, item := range r.Projects {
		go func(item gospeckle.ProjectRequest) {
			p, _, err := c.Project.Create(ctx, item)
			if err != nil {
				ch <- err.Error()
			} else {
				ch <- "Created project " + p.Name + " with ID: " + p.ID
			}
		}(item)
	}

	for _, item := range r.Streams {
		go func(item gospeckle.StreamRequest) {
			s, _, err := c.Stream.Create(ctx, item)
			if err != nil {
				ch <- err.Error()
			} else {
				ch <- "Created stream " + s.Name + " with ID: " + s.StreamID
			}
		}(item)
	}

	for _, item := range r.Clients {
		go func(item gospeckle.APIClientRequest) {
			p, _, err := c.APIClient.Create(ctx, item)
			if err != nil {
				ch <- err.Error()
			} else {
				ch <- "Created client " + p.DocumentName + " with ID: " + p.ID
			}
		}(item)
	}

	for range totalObjects {
		fmt.Println(<-ch)
	}
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

func parseResourceFile(filePath string) (RequestObjects, error) {
	var r RequestObjects

	data, err := os.Open(filePath)
	if err != nil {
		return r, err
	}

	fileType := filepath.Ext(filePath)

	switch fileType {
	case ".yaml":
		d := yaml.NewDecoder(data)
		err := streamDecode(d, &r)
		if err != nil {
			return r, err
		}
	case ".yml":
		d := yaml.NewDecoder(data)
		err := streamDecode(d, &r)
		if err != nil {
			return r, err
		}

	case ".json":
		d := json.NewDecoder(data)
		err := streamDecode(d, &r)
		if err != nil {
			return r, err
		}
	}

	return r, nil
}

func streamDecode(d Decoder, r *RequestObjects) error {

	for {
		var i FileInput
		if err := d.Decode(&i); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		// fmt.Println(i)

		switch i.Type {
		case "project":
			var p gospeckle.ProjectRequest
			err := mapstructure.Decode(i.Spec, &p)
			p.RequestMetadata = i.Metadata
			if err != nil {
				return err
			}
			r.Projects = append(r.Projects, p)

		case "stream":
			var s gospeckle.StreamRequest
			err := mapstructure.Decode(i.Spec, &s)
			s.RequestMetadata = i.Metadata
			if err != nil {
				return err
			}
			r.Streams = append(r.Streams, s)

		case "client":
			var c gospeckle.APIClientRequest
			err := mapstructure.Decode(i.Spec, &c)
			c.RequestMetadata = i.Metadata
			if err != nil {
				return err
			}
			r.Clients = append(r.Clients, c)
		}
	}
	return nil
}

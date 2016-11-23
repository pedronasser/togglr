package main

import (
	"fmt"

	"errors"
	"github.com/jason0x43/go-toggl"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os/user"
	"reflect"
	"strconv"
	"strings"
)

func configCmd() cli.Command {
	return cli.Command{
		Name:      "config",
		Usage:     "togglr configuration",
		ArgsUsage: "name [value]",
		Action:    setGetConfig,
	}
}

var cfg *config

func setGetConfig(c *cli.Context) error {
	if len(c.Args()) == 0 {
		cfgData, err := ioutil.ReadFile(configFile)
		if err != nil {
			return fmt.Errorf("Failed to load config: %s", err)
		}
		fmt.Println(string(cfgData))
		return nil
	}

	name := strings.ToLower(c.Args()[0])

	var value string
	if len(c.Args()) == 2 {
		value = c.Args()[1]
	}

	if value != "" {
		// Set value
		typ := reflect.TypeOf(config{})

		realName := ""
		for i := 0; i < typ.NumField(); i++ {
			field := typ.FieldByIndex([]int{i})
			if strings.ToLower(field.Name) == name {
				realName = field.Name
				break
			}
		}

		field, ok := typ.FieldByName(realName)
		if ok {
			cfgVal := reflect.ValueOf(cfg)
			fieldVal := cfgVal.Elem().FieldByName(realName)
			switch field.Type.Kind() {
			case reflect.Float64:
				v, _ := strconv.ParseFloat(value, 10)
				fieldVal.SetFloat(v)
			case reflect.Int:
				v, _ := strconv.ParseInt(value, 10, 64)
				fieldVal.SetInt(v)
			case reflect.String:
				fieldVal.SetString(value)
			default:
				return errors.New("Invalid config")
			}
			saveConfig()
			fmt.Printf("`%s` is now `%v`\n", name, value)
		} else {
			return fmt.Errorf("Invalid config: %s", name)
		}
	} else {
		// Get value
		cfgData, err := ioutil.ReadFile(configFile)
		if err != nil {
			return fmt.Errorf("Failed to load config: %s", err)
		}
		var c map[string]interface{}
		err = yaml.Unmarshal(cfgData, &c)
		if err != nil {
			return fmt.Errorf("Failed to parse configuration: %s", err)
		}
		fmt.Printf("%s: %v\n", name, c[name])
	}

	return nil
}

type config struct {
	Token     string         `yaml:"token"`
	Workspace int            `yaml:"workspace"`
	Name      string         `yaml:"name"`
	Currency  string         `yaml:"currency"`
	Rate      float64        `yaml:"rate"`
	Aliases   map[int]string `yaml:"aliases"`
}

var configFile = "%s/.togglr.yml"

func init() {
	initConfig()
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("Couldn't locate user's HOME directory: %s", err)
	}
	configFile = fmt.Sprintf(configFile, usr.HomeDir)
	loadConfig()
}

func initConfig() {
	cfg = &config{}
	if cfg.Currency == "" {
		cfg.Currency = "USD"
	}
	cfg.Aliases = make(map[int]string)
}

func loadConfig() error {
	cfgData, err := ioutil.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("Failed to load config: %s", err)
	}
	err = yaml.Unmarshal(cfgData, cfg)
	if err != nil {
		return fmt.Errorf("Failed to parse configuration: %s", err)
	}
	return nil
}

func saveConfig() error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("Failed to marshal configuration: %s", err)
	}
	err = ioutil.WriteFile(configFile, data, 0644)
	if err != nil {
		return fmt.Errorf("Failed to save config: %s", err)
	}
	return nil
}

func doLogin(username, password string) (string, error) {
	token, err := getUserToken(username, password)
	if err != nil {
		return token, fmt.Errorf("Failed to login to toggl.com: %s", err)
	}

	loadConfig()

	cfg.Token = token

	session, err := initSession()
	if err != nil {
		return token, fmt.Errorf("Failed to init session: %s", err)
	}

	acc, err := session.GetAccount()
	if err != nil {
		return token, fmt.Errorf("Failed to get account information: %s", err)
	}

	cfg.Workspace = acc.Data.Workspaces[0].ID

	err = saveConfig()
	if err != nil {
		return token, fmt.Errorf("Failed to save config: %s", err)
	}

	return token, nil
}

func initSession() (toggl.Session, error) {
	if cfg.Token == "" {
		return toggl.Session{}, fmt.Errorf("You are not logged in. Please login to your account: togglr login")
	}

	session := toggl.OpenSession(cfg.Token)
	return session, nil
}

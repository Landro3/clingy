package services

import (
	"clingy-client/util"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Username   string        `json:"username"`
	ServerAddr string        `json:"serverAddr"`
	UniqueID   string        `json:"uniqueId"`
	Contacts   []ContactInfo `json:"contacts"`

	configPath string
}

func NewConfig() *Config {
	configPath := "./config.json"
	exe, err := os.Executable()
	if err == nil {
		configPath = filepath.Join(filepath.Dir(exe), "config.json")
	}
	config := &Config{configPath: configPath}
	err = config.loadFromFile()
	if err != nil {
		util.Log(fmt.Sprintf("error loading config: %s", err))
	}
	return config
}

func (c *Config) UpdateConfig(config *Config) {
	c.ServerAddr = config.ServerAddr
	c.Username = config.Username
	c.UniqueID = config.UniqueID
	err := c.saveToFile()
	if err != nil {
		util.Log(fmt.Sprintf("error saving config: %s", err))
	}
}

func (c *Config) loadFromFile() error {
	err := util.LoadJSONFile(c.configPath, c)
	if err != nil {
		c.Username = ""
		c.ServerAddr = ""
		c.UniqueID = ""
		c.Contacts = []ContactInfo{}
	}

	if c.Contacts == nil {
		c.Contacts = []ContactInfo{}
	}

	return nil
}

func (c *Config) saveToFile() error {
	err := util.SaveToJSONFile(c.configPath, c)
	if err != nil {
		util.Log(fmt.Sprintf("error saving config: %s", err))
	}
	return nil
}

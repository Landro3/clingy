package services

import (
	"clingy-client/util"
	"fmt"
)

type Config struct {
	Username   string        `json:"username"`
	ServerAddr string        `json:"serverAddr"`
	UniqueID   string        `json:"uniqueId"`
	Contacts   []ContactInfo `json:"contacts"`
}

func NewConfig() *Config {
	config := &Config{}
	err := config.loadFromFile()
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
	err := util.LoadJSONFile("./config.json", c)
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
	err := util.SaveToJSONFile("./config.json", c)
	if err != nil {
		util.Log(fmt.Sprintf("error saving config: %s", err))
	}
	return nil
}

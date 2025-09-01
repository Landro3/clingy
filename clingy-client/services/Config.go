package services

import (
	"clingy-client/util"
	"fmt"
)

type Config struct {
	Username   string `json:"username"`
	ServerAddr string `json:"server_addr"`
	UniqueID   string `json:"unique_id"`
}

func NewConfig() *Config {
	config := &Config{}
	config.loadFromFile()
	return config
}

func (c *Config) UpdateConfig(config *Config) {
	c.ServerAddr = config.ServerAddr
	c.Username = config.Username
	c.UniqueID = config.UniqueID
	c.saveToFile()
}

func (c *Config) loadFromFile() error {
	err := util.LoadJSONFile("./config.json", c)
	if err != nil {
		c.Username = ""
		c.ServerAddr = ""
		c.UniqueID = ""
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

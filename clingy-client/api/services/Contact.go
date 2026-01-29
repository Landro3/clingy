package services

import (
	"clingy-client/util"
	"fmt"
)

type ContactInfo struct {
	Username string        `json:"username"`
	ID       string        `json:"uniqueId"`
	Contacts []ContactInfo `json:"contacts"`
}

func (c ContactInfo) Title() string {
	return c.Username
}

func (c ContactInfo) Description() string {
	return c.ID
}

type Contact struct {
	config *Config
}

func NewContact(config *Config) *Contact {
	contact := &Contact{
		config: config,
	}
	return contact
}

func (c *Contact) AddContact(ci ContactInfo) {
	c.config.Contacts = append(c.config.Contacts, ci)
	c.saveToConfig()
}

func (c Contact) saveToConfig() {
	err := c.config.saveToFile()
	if err != nil {
		util.Log(fmt.Sprintf("error saving config: %s", err))
	}
}

func NewContactInfo(username, id string) ContactInfo {
	return ContactInfo{
		Username: username,
		ID:       id,
	}
}

func (c *Contact) RemoveContact(ID string) {
	for i, item := range c.config.Contacts {
		if item.ID == ID {
			c.config.Contacts = append(c.config.Contacts[:i], c.config.Contacts[i+1:]...)
			break
		}
	}
	c.saveToConfig()
}

func (c Contact) UpdateContact(index int, ci ContactInfo) {
	if index >= 0 && index < len(c.config.Contacts) {
		c.config.Contacts[index] = ci
		c.saveToConfig()
	}
}

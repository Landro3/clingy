package services

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
)

type ContactInfo struct {
	Username string `json:"username"`
	ID       string `json:"id"`
}

func (c ContactInfo) FilterValue() string {
	return c.Username
}

func (c ContactInfo) Title() string {
	return c.Username
}

func (c ContactInfo) Description() string {
	return fmt.Sprintf("%s", c.ID)
}

type Contact struct {
	Contacts []ContactInfo
	config   *Config
}

func NewContact(config *Config) *Contact {
	contact := &Contact{
		Contacts: config.Contacts,
		config:   config,
	}
	return contact
}

func (c *Contact) AddContact(ci ContactInfo) {
	c.Contacts = append(c.Contacts, ci)
	c.saveToConfig()
}

func (c *Contact) saveToConfig() {
	c.config.Contacts = c.Contacts
	c.config.saveToFile()
}

func (c *Contact) ToListItems() []list.Item {
	items := make([]list.Item, len(c.Contacts))
	for i, contact := range c.Contacts {
		items[i] = contact
	}
	return items
}

func NewContactInfo(username, id string) ContactInfo {
	return ContactInfo{
		Username: username,
		ID:       id,
	}
}

func (c *Contact) RemoveContact(ID string) {
	for i, item := range c.Contacts {
		if item.ID == ID {
			c.Contacts = append(c.Contacts[:i], c.Contacts[i+1:]...)
			break
		}
	}
	c.saveToConfig()
}

func (c *Contact) UpdateContact(index int, ci ContactInfo) {
	if index >= 0 && index < len(c.Contacts) {
		c.Contacts[index] = ci
		c.saveToConfig()
	}
}

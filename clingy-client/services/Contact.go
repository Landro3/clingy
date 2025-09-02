package services

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
)

type ContactInfo struct {
	Username string
	ID       string
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
}

func NewContact() *Contact {
	contact := &Contact{}
	return contact
}

func (c *Contact) AddContact(ci ContactInfo) {
	c.Contacts = append(c.Contacts, ci)
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

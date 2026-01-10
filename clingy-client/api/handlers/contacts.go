package handlers

import (
	"clingy-client/services"
	"encoding/json"
	"log"
	"net/http"
	"slices"
)

type ContactHandler struct {
	contactService *services.Contact
	configService  *services.Config
}

func NewContactHandler(contactService *services.Contact, configService *services.Config) *ContactHandler {
	return &ContactHandler{
		contactService: contactService,
		configService:  configService,
	}
}

func (h *ContactHandler) GetContacts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	log.Printf("API: GET /api/contacts - Fetching %d contacts", len(h.configService.Contacts))
	contacts := h.configService.Contacts // TODO: read from file every time

	if err := json.NewEncoder(w).Encode(contacts); err != nil {
		log.Printf("API: Error encoding contacts: %v", err)
		http.Error(w, "Failed to encode contacts", http.StatusInternalServerError)
		return
	}
}

func (h *ContactHandler) CreateContact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var contact services.ContactInfo
	if !decodeJSONBody(w, r, &contact) {
		return
	}

	index := slices.IndexFunc(h.configService.Contacts, func(c services.ContactInfo) bool {
		return c.ID == contact.ID
	})
	if index >= 0 {
		log.Printf("API: Contact found with ID: %s", contact.ID)
		http.Error(w, "Contact with UUID already exists", http.StatusConflict)
		return
	}

	log.Printf("API: POST /api/contacts - Creating contact: %s (ID: %s)", contact.Username, contact.ID)
	h.contactService.AddContact(contact)
	log.Printf("API: Contact created successfully. Total contacts: %d", len(h.configService.Contacts))

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(contact); err != nil {
		log.Printf("API: Error encoding created contact: %v", err)
		http.Error(w, "Failed to encode contact", http.StatusInternalServerError)
		return
	}
}

func (h *ContactHandler) UpdateContact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var contact struct {
		services.ContactInfo
		CurrentID string `json:"currentId"`
	}
	if !decodeJSONBody(w, r, &contact) {
		return
	}

	log.Printf("API: PUT /api/contacts - Updating contact: %s (ID: %s)", contact.Username, contact.ID)
	index := slices.IndexFunc(h.configService.Contacts, func(c services.ContactInfo) bool {
		return c.ID == contact.CurrentID
	})

	if index == -1 {
		log.Printf("API: Contact not found with ID: %s", contact.ID)
		http.Error(w, "Contact not found", http.StatusNotFound)
		return
	}

	h.contactService.UpdateContact(index, contact.ContactInfo)
	log.Printf("API: Contact updated successfully")

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(contact); err != nil {
		log.Printf("API: Error encoding updated contact: %v", err)
		http.Error(w, "Failed to encode contact", http.StatusInternalServerError)
		return
	}
}

func (h *ContactHandler) DeleteContact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.URL.Query().Get("id")
	if id == "" {
		log.Printf("API: DELETE /api/contacts - Missing id parameter")
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	log.Printf("API: DELETE /api/contacts - Deleting contact with ID: %s", id)
	h.contactService.RemoveContact(id)
	log.Printf("API: Contact deleted successfully. Total contacts: %d", len(h.configService.Contacts))

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Contact deleted successfully"}`))
}

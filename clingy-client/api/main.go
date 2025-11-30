package main

import (
	"clingy-client/handlers"
	"clingy-client/services"
	"log"
)

func main() {
	log.Println("=== Clingy API Server ===")

	// Initialize services
	log.Println("Initializing services...")
	configService := services.NewConfig()
	contactService := services.NewContact(configService)
	quicService := services.NewQuic(configService)

	log.Printf("Config loaded - Username: %s, ServerAddr: %s",
		configService.Username, configService.ServerAddr)
	log.Printf("Loaded %d existing contacts", len(contactService.Contacts))

	// Start HTTP server
	server := handlers.NewServer(configService, contactService, quicService)

	if err := server.Start(":8888"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

package main

import (
	"clingy-client/handlers"
	"clingy-client/services"
	"flag"
	"log"
)

func main() {
	log.Println("=== Clingy Client Server ===")

	port := flag.String("port", "8888", "server port")
	flag.Parse()

	// Initialize services
	log.Println("Initializing services...")
	configService := services.NewConfig()
	contactService := services.NewContact(configService)
	quicService := services.NewQuic(configService)

	// Start HTTP server
	server := handlers.NewServer(configService, contactService, quicService)

	if err := server.Start(":" + *port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

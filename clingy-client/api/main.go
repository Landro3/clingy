package main

import (
	"clingy-client/handlers"
	"clingy-client/services"
	"clingy-client/util"
	"flag"
	"log"
)

func main() {
	log.Println("=== Clingy Client Server ===")

	port := flag.String("port", "8888", "server port")
	flag.Parse()

	// Initialize services
	log.Println("Initializing services...")
	chatChannelManager := util.NewChannelManager[services.ChatMessage](100) 
	configService := services.NewConfig()
	contactService := services.NewContact(configService)
	http3Service := services.NewHttp3(configService, chatChannelManager.SendMessage)

	// Start HTTP server
	chatChannel, err := chatChannelManager.GetChannel()
	if err != nil {
		return
	}
	server := handlers.NewServer(configService, contactService, http3Service, chatChannel)

	if err := server.Start(":" + *port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

package main

import (
	"log"
	"clingy-server/services"
)

func main() {
	serverType := services.Http2 // TODO: get from cli
	server := services.NewClingyServer(serverType) // TODO: add connection details like port, host
	err := server.ListenAndServeTLS("", "")
	if err != nil {
		log.Fatal(err)
	}
}


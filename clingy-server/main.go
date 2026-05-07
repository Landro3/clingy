package main

import (
	"clingy-server/services"
)

func main() {
	serverType := services.Http3 // TODO: get from cli
	services.StartClingyServer(serverType) // TODO: add connection details like port, host
}


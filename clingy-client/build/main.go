package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
)

//go:embed clingy-api
var apiBin []byte

//go:embed clingy-ui
var uiBin []byte

var defaultPort = "8888"

func main() {
	port := flag.String("port", defaultPort, "API server port")
	flag.Parse()

	log.Println("=== Clingy Client ===")

	dir, err := os.MkdirTemp("", "clingy-*")
	if err != nil {
		log.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	apiPath := filepath.Join(dir, "api")
	uiPath := filepath.Join(dir, "ui")

	if err := os.WriteFile(apiPath, apiBin, 0755); err != nil {
		log.Fatalf("Failed to write API binary: %v", err)
	}
	if err := os.WriteFile(uiPath, uiBin, 0755); err != nil {
		log.Fatalf("Failed to write UI binary: %v", err)
	}

	logPath := "api.log"
	exe, err := os.Executable()
	if err == nil {
		logPath = filepath.Join(filepath.Dir(exe), "api.log")
	}

	apiLogFile, err := os.Create(logPath)
	if err != nil {
		log.Fatalf("Failed to create API log file: %v", err)
	}
	defer apiLogFile.Close()

	log.Println("Starting API server...")
	api := exec.Command(apiPath, "-port", *port)
	api.Stdout = apiLogFile
	api.Stderr = apiLogFile

	if err := api.Start(); err != nil {
		log.Fatalf("Failed to start API: %v", err)
	}
	defer api.Process.Kill()

	os.Setenv("API_URL", fmt.Sprintf("http://localhost:%s/api", *port))

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	log.Println("Starting UI...")
	ui := exec.Command(uiPath)
	ui.Stdin = os.Stdin
	ui.Stdout = os.Stdout
	ui.Stderr = os.Stderr

	if err := ui.Start(); err != nil {
		api.Process.Kill()
		log.Fatalf("Failed to start UI: %v", err)
	}

	<-sigChan
	log.Println("\nShutting down...")
	log.Println("Goodbye!")
}

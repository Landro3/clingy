package main

import (
	_ "embed"
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

func main() {
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

	log.Println("Starting API server...")
	api := exec.Command(apiPath, "-port", "8888")
	api.Stdout = os.Stdout
	api.Stderr = os.Stderr

	if err := api.Start(); err != nil {
		log.Fatalf("Failed to start API: %v", err)
	}
	defer api.Process.Kill()

	os.Setenv("API_URL", "http://localhost:8888")

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

	go func() {
		<-sigChan
		log.Println("\nShutting down...")
		ui.Process.Kill()
	}()

	ui.Wait()
	log.Println("Goodbye!")
}

package handlers

import (
	"clingy-client/services"
	"encoding/json"
	"log"
	"net/http"
)

type ConfigHandler struct {
	configService *services.Config
	quicService   *services.Quic
}

func NewConfigHandler(
	configService *services.Config,
	quicService *services.Quic,
) *ConfigHandler {
	return &ConfigHandler{
		configService: configService,
		quicService:   quicService,
	}
}

func (h *ConfigHandler) GetServerConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	log.Printf("API: GET /api/config/server - Fetching server config")
	config := services.NewConfig()
	serverConfig := struct {
		Username   string `json:"username"`
		ServerAddr string `json:"serverAddr"`
		UniqueID   string `json:"uniqueId"`
	}{
		Username:   config.Username,
		ServerAddr: config.ServerAddr,
		UniqueID:   config.UniqueID,
	}

	if err := json.NewEncoder(w).Encode(serverConfig); err != nil {
		log.Printf("API: Error encoding server config: %v", err)
		http.Error(w, "Failed to encode server config", http.StatusInternalServerError)
		return
	}
}

func (h *ConfigHandler) SetServerConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var serverConfig struct {
		Username   string `json:"username"`
		ServerAddr string `json:"serverAddr"`
	}

	if !decodeJSONBody(w, r, &serverConfig) {
		return
	}

	log.Printf("API: POST /api/config/server - Setting server config: %s@%s",
		serverConfig.Username, serverConfig.ServerAddr)

	config := h.configService
	config.Username = serverConfig.Username
	config.ServerAddr = serverConfig.ServerAddr

	assignedUUID, err := h.quicService.Register(config.Username)
	log.Printf("Assigned UUID: %v", assignedUUID)
	if err != nil {
		log.Printf("API: Failed to register with server: %v", err)
		http.Error(w, "Failed to register with server", http.StatusInternalServerError)
		return
	}

	config.UniqueID = assignedUUID
	config.UpdateConfig(config)

	log.Printf("API: Registration successful. UUID: %s", assignedUUID)

	// Prepare response including the assigned UUID
	response := struct {
		Username   string `json:"username"`
		ServerAddr string `json:"serverAddr"`
		UniqueID   string `json:"uniqueId"`
	}{
		Username:   serverConfig.Username,
		ServerAddr: serverConfig.ServerAddr,
		UniqueID:   assignedUUID,
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("API: Error encoding server config response: %v", err)
		http.Error(w, "Failed to encode server config", http.StatusInternalServerError)
		return
	}
}

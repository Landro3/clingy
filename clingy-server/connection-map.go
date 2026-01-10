package main

import (
	"log"
	"net/http"
	"sync"
)

type ConnectionMap struct {
	connections map[string]http.ResponseWriter
	mutex       sync.RWMutex
}

func NewConnectionMap() *ConnectionMap {
	return &ConnectionMap{
		connections: make(map[string]http.ResponseWriter),
	}
}

func (cm *ConnectionMap) Add(userID string, conn http.ResponseWriter) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.connections[userID] = conn
}

func (cm *ConnectionMap) Get(userID string) (http.ResponseWriter, bool) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	conn, exists := cm.connections[userID]
	return conn, exists
}

func (cm *ConnectionMap) Remove(userID string) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	delete(cm.connections, userID)
}

func (cm *ConnectionMap) LogConnections() {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	log.Printf("Active connections (%d):", len(cm.connections))
	for userID := range cm.connections {
		log.Printf("  %s -> %s", userID, "connected")
	}
}

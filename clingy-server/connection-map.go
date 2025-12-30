package main

import (
	"log"
	"sync"

	quic "github.com/quic-go/quic-go"
)

type ConnectionMap struct {
	connections map[string]*quic.Conn
	mutex       sync.RWMutex
}

func NewConnectionMap() *ConnectionMap {
	return &ConnectionMap{
		connections: make(map[string]*quic.Conn),
	}
}

func (cm *ConnectionMap) Add(userID string, conn *quic.Conn) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.connections[userID] = conn
}

func (cm *ConnectionMap) Get(userID string) (*quic.Conn, bool) {
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
	for userID, conn := range cm.connections {
		log.Printf("  %s -> %s", userID, conn.RemoteAddr())
	}
}

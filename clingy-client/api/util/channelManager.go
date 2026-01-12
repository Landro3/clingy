package util

import (
	"errors"
	"fmt"
)

type ChannelManager[T any] struct {
	hasListener bool
	channel chan T
}

func NewChannelManager[T any](bufferSize int) *ChannelManager[T] {
	return &ChannelManager[T]{
		channel: make(chan T, bufferSize),
	}
}

func (cm *ChannelManager[T]) GetChannel() (<-chan T, error) {
	if cm.hasListener {
		return nil, errors.New("listener already registered")
	}

	cm.hasListener = true
	return cm.channel, nil
}

func (cm *ChannelManager[T]) SendMessage(message T) {
	select {
	case cm.channel <- message:
		Log(fmt.Sprintf("📨 Pushing message to channel:\n%v", message))
	default:
		Log("Message channel full, dropping message")
	}
}


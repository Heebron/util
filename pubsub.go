package util

import (
	"sync"
)

type PubSub[T any] struct {
	sync.RWMutex
	clients []chan T
}

// NewPubSub creates a new instance of the PubSub struct with the specified type.
func NewPubSub[T any]() *PubSub[T] {
	return &PubSub[T]{}
}

// Register registers a new reception channel of size buffSize. If buffSize == 0 then broadcasting to this channel is synchronous.
func (s *PubSub[T]) Register(buffSize uint) <-chan T {
	c := make(chan T, buffSize)
	s.Lock()
	s.clients = append(s.clients, c)
	s.Unlock()
	return c
}

// Broadcast sends a message to all registered channels in the PubSub.
// It acquires a read-lock to prevent changes to the client's slice during the broadcast.
// The client's slice is copied into a new buffer for non-concurrent iteration.
// Sends the message to each channel in the buffer.
// Since the channels are buffered, the broadcast is non-blocking.
// Once the broadcast is complete, the read-lock is released.
func (s *PubSub[T]) Broadcast(msg T) {
	s.RLock()
	buffer := make([]chan T, len(s.clients))
	copy(buffer, s.clients)
	s.RUnlock()
	for _, v := range buffer {
		v <- msg
	}
}

// Unregister removes the given reception channel from the PubSub. It closes the channel and removes it from the list of clients.
// If the channel is not found, nothing happens.
func (s *PubSub[T]) Unregister(c <-chan T) {
	s.Lock()
	for i, v := range s.clients {
		if v == c {
			close(v)
			// overwrite and chop method (most efficient use of memory)
			s.clients[i] = s.clients[len(s.clients)-1] // copy the last to the removed
			s.clients = s.clients[:len(s.clients)-1]   // shorten the slice
			break
		}
	}
	s.Unlock()
}

// Size returns the number of registered reception channels in the PubSub struct.
func (s *PubSub[T]) Size() int {
	return len(s.clients)
}

// IsActive returns true if there are active clients subscribed to the PubSub object, otherwise false.
func (s *PubSub[T]) IsActive() bool {
	return len(s.clients) > 0
}

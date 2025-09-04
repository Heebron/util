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
// It acquires a read-lock to prevent changes to the clients slice during the snapshot.
// The clients slice is copied into a new buffer for non-concurrent iteration, then the lock is released.
// Sends the message to each channel in the buffer.
// Note: send may block if a channel has insufficient buffer or no receiver.
func (s *PubSub[T]) Broadcast(msg T) {
	s.RLock()
	buffer := make([]chan T, len(s.clients))
	copy(buffer, s.clients)
	s.RUnlock()
	for _, v := range buffer {
		v <- msg
	}
}

// Unregister removes the given reception channel from the PubSub. It removes it from the list of clients.
// If the channel is not found, nothing happens. The channel is not closed by PubSub to avoid races with concurrent senders.
func (s *PubSub[T]) Unregister(c <-chan T) {
	s.Lock()
	defer s.Unlock()

	for i, v := range s.clients {
		if v == c {
			// overwrite and chop method (most efficient use of memory)
			s.clients[i] = s.clients[len(s.clients)-1] // copy the last to the removed
			s.clients = s.clients[:len(s.clients)-1]   // shorten the slice
			break
		}
	}
}

// Size returns the number of registered reception channels in the PubSub struct.
func (s *PubSub[T]) Size() int {
	s.RLock()
	defer s.RUnlock()
	return len(s.clients)
}

// IsActive returns true if there are active clients subscribed to the PubSub object, otherwise false.
func (s *PubSub[T]) IsActive() bool {
	s.RLock()
	defer s.RUnlock()
	return len(s.clients) > 0
}

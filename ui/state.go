package ui

import (
	"sync"
)

// Invalidator is implemented by windows and other render hosts that can
// schedule a new frame.
type Invalidator interface {
	Invalidate()
}

// State is a reactive value that triggers redraw when changed
type State[T comparable] struct {
	value       T
	mu          sync.RWMutex
	subscribers map[uint64]Invalidator
	nextID      uint64
}

// NewState creates a new reactive state
func NewState[T comparable](initial T) *State[T] {
	return &State[T]{value: initial}
}

// Get the current value (safe for reading)
func (s *State[T]) Get() T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.value
}

// Set a new value and trigger redraw if window is attached
func (s *State[T]) Set(val T) {
	s.mu.Lock()
	changed := s.value != val
	s.value = val
	s.mu.Unlock()

	s.invalidate(changed)
}

// Update applies a function to the current value and sets the result
func (s *State[T]) Update(fn func(T) T) {
	s.mu.Lock()
	oldVal := s.value
	newVal := fn(oldVal)
	changed := oldVal != newVal
	s.value = newVal
	s.mu.Unlock()

	s.invalidate(changed)
}

// Bind sets up the state for reactivity in a render host.
// Use Subscribe when the binding needs an explicit lifecycle.
func (s *State[T]) Bind(invalidator Invalidator) *State[T] {
	s.Subscribe(invalidator)
	return s
}

// Subscribe invalidates a render host whenever the value changes. The
// returned function removes the subscription and is safe to call repeatedly.
func (s *State[T]) Subscribe(invalidator Invalidator) func() {
	if invalidator == nil {
		return func() {}
	}

	s.mu.Lock()
	if s.subscribers == nil {
		s.subscribers = make(map[uint64]Invalidator)
	}
	id := s.nextID
	s.nextID++
	s.subscribers[id] = invalidator
	s.mu.Unlock()

	var once sync.Once
	return func() {
		once.Do(func() {
			s.mu.Lock()
			delete(s.subscribers, id)
			s.mu.Unlock()
		})
	}
}

func (s *State[T]) invalidate(changed bool) {
	if !changed {
		return
	}

	s.mu.RLock()
	subscribers := make([]Invalidator, 0, len(s.subscribers))
	for _, subscriber := range s.subscribers {
		subscribers = append(subscribers, subscriber)
	}
	s.mu.RUnlock()

	for _, subscriber := range subscribers {
		subscriber.Invalidate()
	}
}

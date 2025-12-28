package ui

import (
	"sync"

	"gioui.org/app"
)

// State is a reactive value that triggers redraw when changed
type State[T comparable] struct {
	value  T
	mu     sync.RWMutex
	window *app.Window
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

	if changed && s.window != nil {
		s.window.Invalidate()
	}
}

// Update applies a function to the current value and sets the result
func (s *State[T]) Update(fn func(T) T) {
	s.mu.Lock()
	oldVal := s.value
	newVal := fn(oldVal)
	changed := oldVal != newVal
	s.value = newVal
	s.mu.Unlock()

	if changed && s.window != nil {
		s.window.Invalidate()
	}
}

// Bind sets up the state for reactivity in this app window
func (s *State[T]) Bind(w *app.Window) *State[T] {
	s.window = w
	return s
}

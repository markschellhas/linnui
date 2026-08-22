package ui

import (
	"sync/atomic"
	"testing"
)

type testInvalidator struct {
	count atomic.Int32
}

func (invalidator *testInvalidator) Invalidate() {
	invalidator.count.Add(1)
}

func TestStateSubscriptions(t *testing.T) {
	state := NewState(1)
	first := new(testInvalidator)
	second := new(testInvalidator)
	unsubscribe := state.Subscribe(first)
	state.Subscribe(second)

	state.Set(2)
	state.Set(2)
	if got := first.count.Load(); got != 1 {
		t.Fatalf("first invalidations = %d, want 1", got)
	}
	if got := second.count.Load(); got != 1 {
		t.Fatalf("second invalidations = %d, want 1", got)
	}

	unsubscribe()
	unsubscribe()
	state.Update(func(value int) int { return value + 1 })
	if got := first.count.Load(); got != 1 {
		t.Fatalf("unsubscribed invalidations = %d, want 1", got)
	}
	if got := second.count.Load(); got != 2 {
		t.Fatalf("remaining invalidations = %d, want 2", got)
	}
}

func TestStateBindReturnsState(t *testing.T) {
	state := NewState("initial")
	invalidator := new(testInvalidator)
	if got := state.Bind(invalidator); got != state {
		t.Fatal("Bind must return the same State")
	}
	state.Set("updated")
	if got := invalidator.count.Load(); got != 1 {
		t.Fatalf("invalidations = %d, want 1", got)
	}
}

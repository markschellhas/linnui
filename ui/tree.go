package ui

import (
	"strings"
	"sync"

	"gioui.org/widget"
)

type widgetKind uint8

const (
	buttonWidget widgetKind = iota
	textFieldWidget
	scrollViewWidget
	listViewWidget
	checkboxWidget
	switchWidget
	radioWidget
	sliderWidget
	cardWidget
)

type widgetKey struct {
	kind  widgetKind
	scope string
	id    string
}

type treeStore struct {
	mu      sync.Mutex
	entries map[widgetKey]any
	closed  bool
}

// Tree owns the persistent interaction state for one UI tree or window.
// Create it once and reuse it when rebuilding declarative widgets.
type Tree struct {
	store  *treeStore
	scope  string
	root   bool
	closed bool
}

// NewTree creates an isolated widget-state tree.
func NewTree(_ Invalidator) *Tree {
	return &Tree{
		store: &treeStore{entries: make(map[widgetKey]any)},
		root:  true,
	}
}

// Scope creates a child namespace. IDs can be reused safely across scopes.
func (t *Tree) Scope(id string) *Tree {
	if id == "" {
		panic("ui: tree scope ID must not be empty")
	}
	t.ensureOpen()
	scope := id
	if t.scope != "" {
		scope = t.scope + "\x00" + id
	}
	return &Tree{store: t.store, scope: scope}
}

// Delete removes every kind of widget state with id in this scope.
func (t *Tree) Delete(id string) {
	t.ensureOpen()
	t.store.mu.Lock()
	defer t.store.mu.Unlock()
	for key := range t.store.entries {
		if key.scope == t.scope && key.id == id {
			delete(t.store.entries, key)
		}
	}
}

// Reset removes all widget state in this scope and its descendants.
func (t *Tree) Reset() {
	t.ensureOpen()
	t.store.mu.Lock()
	defer t.store.mu.Unlock()
	t.deleteScopeLocked()
}

// Close releases this scope. Closing the root releases the whole tree.
func (t *Tree) Close() {
	if t == nil || t.closed {
		return
	}
	t.store.mu.Lock()
	defer t.store.mu.Unlock()
	if t.root {
		clear(t.store.entries)
		t.store.closed = true
	} else if !t.store.closed {
		t.deleteScopeLocked()
	}
	t.closed = true
}

func (t *Tree) deleteScopeLocked() {
	for key := range t.store.entries {
		if key.scope == t.scope || strings.HasPrefix(key.scope, t.scope+"\x00") {
			delete(t.store.entries, key)
		}
	}
}

func (t *Tree) ensureOpen() {
	if t == nil {
		panic("ui: nil Tree")
	}
	if t.closed {
		panic("ui: Tree scope is closed")
	}
	t.store.mu.Lock()
	closed := t.store.closed
	t.store.mu.Unlock()
	if closed {
		panic("ui: Tree is closed")
	}
}

func (t *Tree) state(kind widgetKind, id string, create func() any) any {
	t.ensureOpen()
	key := widgetKey{kind: kind, scope: t.scope, id: id}
	t.store.mu.Lock()
	defer t.store.mu.Unlock()
	if value, ok := t.store.entries[key]; ok {
		return value
	}
	value := create()
	t.store.entries[key] = value
	return value
}

type textFieldState struct {
	mu       sync.Mutex
	editor   widget.Editor
	snapshot string
	pending  *string
}

func (t *Tree) clickable(kind widgetKind, id string) *widget.Clickable {
	return t.state(kind, id, func() any { return new(widget.Clickable) }).(*widget.Clickable)
}

func (t *Tree) textField(id string) *textFieldState {
	return t.state(textFieldWidget, id, func() any {
		state := new(textFieldState)
		state.editor.SingleLine = true
		return state
	}).(*textFieldState)
}

func (t *Tree) list(kind widgetKind, id string) *widget.List {
	return t.state(kind, id, func() any { return new(widget.List) }).(*widget.List)
}

func (t *Tree) boolControl(kind widgetKind, id string) *widget.Bool {
	return t.state(kind, id, func() any { return new(widget.Bool) }).(*widget.Bool)
}

func (t *Tree) enum(id string) *widget.Enum {
	return t.state(radioWidget, id, func() any { return new(widget.Enum) }).(*widget.Enum)
}

func (t *Tree) float(id string) *widget.Float {
	return t.state(sliderWidget, id, func() any { return new(widget.Float) }).(*widget.Float)
}

// TextFieldValue returns the latest laid-out value for a text field.
func (t *Tree) TextFieldValue(id string) (string, bool) {
	t.ensureOpen()
	key := widgetKey{kind: textFieldWidget, scope: t.scope, id: id}
	t.store.mu.Lock()
	defer t.store.mu.Unlock()
	state, ok := t.store.entries[key].(*textFieldState)
	if !ok {
		return "", false
	}
	state.mu.Lock()
	defer state.mu.Unlock()
	return state.snapshot, true
}

// SetTextFieldValue queues a value for the next text-field layout.
func (t *Tree) SetTextFieldValue(id, text string) bool {
	t.ensureOpen()
	key := widgetKey{kind: textFieldWidget, scope: t.scope, id: id}
	t.store.mu.Lock()
	defer t.store.mu.Unlock()
	state, ok := t.store.entries[key].(*textFieldState)
	if !ok {
		return false
	}
	state.mu.Lock()
	defer state.mu.Unlock()
	state.pending = &text
	state.snapshot = text
	return true
}

var legacyTree = NewTree(nil)

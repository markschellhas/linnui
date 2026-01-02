package ui

import (
	"sync"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// ScrollDirection defines the scroll axis
type ScrollDirection int

const (
	// ScrollVertical allows vertical scrolling (default)
	ScrollVertical ScrollDirection = iota
	// ScrollHorizontal allows horizontal scrolling
	ScrollHorizontal
)

// ScrollViewOption configures the ScrollView
type ScrollViewOption func(*scrollViewModel)

// Direction sets the scroll direction
func Direction(d ScrollDirection) ScrollViewOption {
	return func(s *scrollViewModel) { s.direction = d }
}

// ScrollID sets a unique ID for the scroll view (for state persistence)
// Use this when you have multiple scroll views
func ScrollID(id string) ScrollViewOption {
	return func(s *scrollViewModel) { s.id = id }
}

// scrollViewModel holds ScrollView configuration (internal)
type scrollViewModel struct {
	id        string
	direction ScrollDirection
	child     Widget
}

// scrollRegistry stores list state by ID
var (
	scrollRegistry = make(map[string]*widget.List)
	scrollMu       sync.Mutex
)

// getList returns a persistent list for the given ID
func getList(id string) *widget.List {
	scrollMu.Lock()
	defer scrollMu.Unlock()

	if l, ok := scrollRegistry[id]; ok {
		return l
	}
	l := new(widget.List)
	scrollRegistry[id] = l
	return l
}

// ScrollView creates a scrollable container for a single child
// Usage: ScrollView(child, Direction(ScrollVertical))
func ScrollView(child Widget, opts ...ScrollViewOption) Widget {
	s := &scrollViewModel{
		id:        "default-scroll", // Default ID
		direction: ScrollVertical,   // sensible default
		child:     child,
	}
	for _, opt := range opts {
		opt(s)
	}

	// Get persistent list state
	list := getList(s.id)

	// Set the axis based on direction
	if s.direction == ScrollHorizontal {
		list.Axis = layout.Horizontal
	} else {
		list.Axis = layout.Vertical
	}

	return func(gtx layout.Context, th *Theme) layout.Dimensions {
		if s.child == nil {
			return layout.Dimensions{}
		}

		child := s.child
		return material.List(th.Theme, list).Layout(gtx, 1, func(gtx layout.Context, _ int) layout.Dimensions {
			return child(gtx, th)
		})
	}
}

// ListView creates a scrollable list of children with efficient rendering
// Unlike ScrollView which wraps a single child, ListView is optimized for many items
// Usage: ListView([]Widget{item1, item2, ...}, Direction(ScrollVertical))
func ListView(children []Widget, opts ...ScrollViewOption) Widget {
	s := &scrollViewModel{
		id:        "default-listview",
		direction: ScrollVertical,
	}
	for _, opt := range opts {
		opt(s)
	}

	// Get persistent list state
	list := getList(s.id)

	// Set the axis based on direction
	if s.direction == ScrollHorizontal {
		list.Axis = layout.Horizontal
	} else {
		list.Axis = layout.Vertical
	}

	return func(gtx layout.Context, th *Theme) layout.Dimensions {
		return material.List(th.Theme, list).Layout(gtx, len(children), func(gtx layout.Context, i int) layout.Dimensions {
			if i < len(children) && children[i] != nil {
				return children[i](gtx, th)
			}
			return layout.Dimensions{}
		})
	}
}

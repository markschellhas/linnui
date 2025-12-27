// Package ui provides a Flutter/SwiftUI-like declarative developer experience for Go.
// All widgets and layouts are accessible via the ui. prefix.
package ui

import (
	"gioui.org/layout"
)

// Widget is the core type: a function that lays out UI
type Widget func(gtx layout.Context, th *Theme) layout.Dimensions

// MainAxisAlignment controls alignment along the main axis
type MainAxisAlignment int

const (
	MainAxisStart MainAxisAlignment = iota
	MainAxisCenter
	MainAxisEnd
	MainAxisSpaceBetween
	MainAxisSpaceAround
	MainAxisSpaceEvenly
)

// CrossAxisAlignment controls alignment along the cross axis
type CrossAxisAlignment int

const (
	CrossAxisStart CrossAxisAlignment = iota
	CrossAxisCenter
	CrossAxisEnd
	CrossAxisStretch
	CrossAxisBaseline
)

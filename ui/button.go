package ui

import (
	"image"
	"image/color"
	"sync"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// buttonRegistry stores clickable state by button label
var (
	buttonRegistry = make(map[string]*widget.Clickable)
	buttonMu       sync.Mutex
)

// ButtonVariant defines the visual style of a button
type ButtonVariant int

const (
	Filled ButtonVariant = iota
	Outlined
	TextButton
	Elevated
)

// ButtonOption configures the Button
type ButtonOption func(*buttonModel)

// OnClick sets the click handler for a Button
func OnClick(fn func()) ButtonOption {
	return func(b *buttonModel) { b.onClick = fn }
}

// Variant sets the button style variant
func Variant(v ButtonVariant) ButtonOption {
	return func(b *buttonModel) { b.variant = v }
}

// ButtonID sets a unique ID for the button (for state persistence)
// Use this when you have multiple buttons with the same label
func ButtonID(id string) ButtonOption {
	return func(b *buttonModel) { b.id = id }
}

// buttonModel holds button state and configuration (internal)
type buttonModel struct {
	id      string
	label   string
	variant ButtonVariant
	onClick func()
}

// getClickable returns a persistent clickable for the given ID
func getClickable(id string) *widget.Clickable {
	buttonMu.Lock()
	defer buttonMu.Unlock()

	if c, ok := buttonRegistry[id]; ok {
		return c
	}
	c := new(widget.Clickable)
	buttonRegistry[id] = c
	return c
}

// Button creates a clickable button widget
func Button(label string, opts ...ButtonOption) Widget {
	b := &buttonModel{
		id:      label, // Default ID is the label
		label:   label,
		variant: Filled, // sensible default
	}
	for _, opt := range opts {
		opt(b)
	}

	// Get persistent clickable using the ID
	clickable := getClickable(b.id)
	onClick := b.onClick // Capture the handler

	return func(gtx layout.Context, th *Theme) layout.Dimensions {
		// Handle clicks
		for clickable.Clicked(gtx) {
			if onClick != nil {
				onClick()
			}
		}

		// Base material button with ripple
		mat := material.Button(th.Theme, clickable, b.label)

		// Apply variant-specific styling
		switch b.variant {
		case Filled:
			mat.Background = th.Palette.Primary
			mat.Color = th.Palette.OnPrimary
			mat.CornerRadius = unit.Dp(12)
		case Outlined:
			mat.Background = color.NRGBA{A: 0} // Transparent
			mat.Color = th.Palette.Primary
			mat.CornerRadius = unit.Dp(12)
			// Draw button first, then add outline
			return layout.Stack{}.Layout(gtx,
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return mat.Layout(gtx)
				}),
				layout.Expanded(func(gtx layout.Context) layout.Dimensions {
					size := gtx.Constraints.Min
					radius := gtx.Dp(unit.Dp(12))
					borderWidth := gtx.Dp(unit.Dp(1))

					// Draw outline using stroke
					rect := image.Rect(0, 0, size.X, size.Y)
					outline := clip.Stroke{
						Path:  clip.UniformRRect(rect, radius).Path(gtx.Ops),
						Width: float32(borderWidth),
					}.Op().Push(gtx.Ops)
					paint.Fill(gtx.Ops, th.Palette.Primary)
					outline.Pop()

					return layout.Dimensions{Size: size}
				}),
			)
		case TextButton:
			mat.Background = color.NRGBA{A: 0} // Transparent
			mat.Color = th.Palette.Primary
			mat.CornerRadius = unit.Dp(12)
		case Elevated:
			mat.Background = th.Palette.SurfaceVariant
			mat.Color = th.Palette.Primary
			mat.CornerRadius = unit.Dp(12)
		}

		return mat.Layout(gtx)
	}
}

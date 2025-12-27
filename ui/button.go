package ui

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
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

// buttonModel holds button state and configuration (internal)
type buttonModel struct {
	clickable *widget.Clickable
	label     string
	variant   ButtonVariant
	onClick   func()
}

// Button creates a clickable button widget
func Button(label string, opts ...ButtonOption) Widget {
	b := &buttonModel{
		clickable: new(widget.Clickable),
		label:     label,
		variant:   Filled, // sensible default
	}
	for _, opt := range opts {
		opt(b)
	}

	return func(gtx layout.Context, th *Theme) layout.Dimensions {
		// Handle clicks
		for b.clickable.Clicked(gtx) {
			if b.onClick != nil {
				b.onClick()
			}
		}

		// Base material button with ripple
		mat := material.Button(th.Theme, b.clickable, b.label)

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

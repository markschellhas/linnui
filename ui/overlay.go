package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
)

// Overlay is a widget layer rendered above Scaffold content.
type Overlay struct {
	widget  Widget
	visible func() bool
	modal   bool
}

// OverlayOption configures a custom Overlay.
type OverlayOption func(*Overlay)

// OverlayVisible displays an overlay only while visible returns true.
func OverlayVisible(visible func() bool) OverlayOption {
	return func(overlay *Overlay) { overlay.visible = visible }
}

// OverlayModal prevents underlying Scaffold content from receiving input.
func OverlayModal() OverlayOption {
	return func(overlay *Overlay) { overlay.modal = true }
}

// CustomOverlay creates an arbitrary Scaffold overlay.
func CustomOverlay(widget Widget, opts ...OverlayOption) Overlay {
	overlay := Overlay{widget: widget}
	for _, opt := range opts {
		opt(&overlay)
	}
	return overlay
}

func (overlay Overlay) active() bool {
	return overlay.widget != nil && (overlay.visible == nil || overlay.visible())
}

func alignedOverlay(direction layout.Direction, insets Insets, widget Widget) Widget {
	return func(gtx layout.Context, th *Theme) layout.Dimensions {
		return direction.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Top:    unit.Dp(insets.Top),
				Right:  unit.Dp(insets.Right),
				Bottom: unit.Dp(insets.Bottom),
				Left:   unit.Dp(insets.Left),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return widget(gtx, th)
			})
		})
	}
}

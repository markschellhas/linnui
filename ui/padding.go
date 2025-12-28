package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
)

// Insets defines padding/margin values
type Insets struct {
	Top    float32
	Right  float32
	Bottom float32
	Left   float32
}

// InsetsAll creates uniform insets on all sides
func InsetsAll(dp float32) Insets {
	return Insets{Top: dp, Right: dp, Bottom: dp, Left: dp}
}

// InsetsSymmetric creates symmetric horizontal and vertical insets
func InsetsSymmetric(horizontal, vertical float32) Insets {
	return Insets{Top: vertical, Right: horizontal, Bottom: vertical, Left: horizontal}
}

// InsetsOnly creates insets with specific sides
func InsetsOnly(top, right, bottom, left float32) Insets {
	return Insets{Top: top, Right: right, Bottom: bottom, Left: left}
}

// Padding creates space inside a widget around its child
// Usage: Padding(InsetsAll(16), child)
func Padding(insets Insets, child Widget) Widget {
	return func(gtx layout.Context, th *Theme) layout.Dimensions {
		return layout.Inset{
			Top:    unit.Dp(insets.Top),
			Right:  unit.Dp(insets.Right),
			Bottom: unit.Dp(insets.Bottom),
			Left:   unit.Dp(insets.Left),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			if child != nil {
				return child(gtx, th)
			}
			return layout.Dimensions{}
		})
	}
}

// Margin creates space outside a widget (alias for Padding in Gio's model)
// Usage: Margin(InsetsAll(8), child)
func Margin(insets Insets, child Widget) Widget {
	return Padding(insets, child)
}

package ui

import (
	"image"

	"gioui.org/layout"
	"gioui.org/unit"
)

// SizedBoxOption configures the SizedBox
type SizedBoxOption func(*sizedBoxModel)

// Width sets the fixed width in dp
func Width(dp float32) SizedBoxOption {
	return func(s *sizedBoxModel) { s.width = dp; s.hasWidth = true }
}

// Height sets the fixed height in dp
func Height(dp float32) SizedBoxOption {
	return func(s *sizedBoxModel) { s.height = dp; s.hasHeight = true }
}

// sizedBoxModel holds SizedBox configuration (internal)
type sizedBoxModel struct {
	width     float32
	height    float32
	hasWidth  bool
	hasHeight bool
}

// SizedBox creates a box with fixed dimensions
// Usage: SizedBox(Width(100), Height(50), child) or SizedBox(Height(20)) for spacing
func SizedBox(opts ...any) Widget {
	s := &sizedBoxModel{}
	var child Widget

	for _, opt := range opts {
		switch v := opt.(type) {
		case SizedBoxOption:
			v(s)
		case Widget:
			child = v
		}
	}

	return func(gtx layout.Context, th *Theme) layout.Dimensions {
		// Calculate constraints
		minW, maxW := 0, gtx.Constraints.Max.X
		minH, maxH := 0, gtx.Constraints.Max.Y

		if s.hasWidth {
			w := gtx.Dp(unit.Dp(s.width))
			minW, maxW = w, w
		}
		if s.hasHeight {
			h := gtx.Dp(unit.Dp(s.height))
			minH, maxH = h, h
		}

		// Apply constraints to child
		childGtx := gtx
		childGtx.Constraints = layout.Constraints{
			Min: image.Pt(minW, minH),
			Max: image.Pt(maxW, maxH),
		}

		if child != nil {
			return child(childGtx, th)
		}

		// No child - just return the sized space
		return layout.Dimensions{
			Size: image.Pt(
				gtx.Constraints.Constrain(image.Pt(maxW, maxH)).X,
				gtx.Constraints.Constrain(image.Pt(maxW, maxH)).Y,
			),
		}
	}
}

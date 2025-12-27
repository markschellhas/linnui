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

// Child sets the child widget inside the SizedBox
func Child(w Widget) SizedBoxOption {
	return func(s *sizedBoxModel) { s.child = w }
}

// sizedBoxModel holds SizedBox configuration (internal)
type sizedBoxModel struct {
	width     float32
	height    float32
	hasWidth  bool
	hasHeight bool
	child     Widget
}

// SizedBox creates a box with fixed dimensions
// Use Width(), Height(), and Child() options to configure
func SizedBox(opts ...SizedBoxOption) Widget {
	s := &sizedBoxModel{}
	for _, opt := range opts {
		opt(s)
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

		if s.child != nil {
			return s.child(childGtx, th)
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

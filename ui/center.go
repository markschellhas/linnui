package ui

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op"
)

// Center centers its child within the available space
// Usage: Center(child)
func Center(child Widget) Widget {
	return func(gtx layout.Context, th *Theme) layout.Dimensions {
		// First, we need to measure the child to know its size
		macro := op.Record(gtx.Ops)
		childDims := layout.Dimensions{}
		if child != nil {
			childDims = child(gtx, th)
		}
		call := macro.Stop()

		// Calculate centered position
		availableSize := gtx.Constraints.Max
		offsetX := (availableSize.X - childDims.Size.X) / 2
		offsetY := (availableSize.Y - childDims.Size.Y) / 2

		if offsetX < 0 {
			offsetX = 0
		}
		if offsetY < 0 {
			offsetY = 0
		}

		// Apply offset and draw
		defer op.Offset(image.Pt(offsetX, offsetY)).Push(gtx.Ops).Pop()
		call.Add(gtx.Ops)

		return layout.Dimensions{
			Size: availableSize,
		}
	}
}

// Align positions its child according to the specified alignment
// Usage: Align(TopLeft, child)
func Align(alignment Alignment, child Widget) Widget {
	return func(gtx layout.Context, th *Theme) layout.Dimensions {
		// Measure the child first
		macro := op.Record(gtx.Ops)
		childDims := layout.Dimensions{}
		if child != nil {
			childDims = child(gtx, th)
		}
		call := macro.Stop()

		availableSize := gtx.Constraints.Max
		var offsetX, offsetY int

		// Calculate X offset
		switch alignment {
		case TopLeft, CenterLeft, BottomLeft:
			offsetX = 0
		case TopCenter, CenterCenter, BottomCenter:
			offsetX = (availableSize.X - childDims.Size.X) / 2
		case TopRight, CenterRight, BottomRight:
			offsetX = availableSize.X - childDims.Size.X
		}

		// Calculate Y offset
		switch alignment {
		case TopLeft, TopCenter, TopRight:
			offsetY = 0
		case CenterLeft, CenterCenter, CenterRight:
			offsetY = (availableSize.Y - childDims.Size.Y) / 2
		case BottomLeft, BottomCenter, BottomRight:
			offsetY = availableSize.Y - childDims.Size.Y
		}

		if offsetX < 0 {
			offsetX = 0
		}
		if offsetY < 0 {
			offsetY = 0
		}

		// Apply offset and draw
		defer op.Offset(image.Pt(offsetX, offsetY)).Push(gtx.Ops).Pop()
		call.Add(gtx.Ops)

		return layout.Dimensions{
			Size: availableSize,
		}
	}
}

// Alignment defines positioning within a container
type Alignment int

const (
	TopLeft Alignment = iota
	TopCenter
	TopRight
	CenterLeft
	CenterCenter
	CenterRight
	BottomLeft
	BottomCenter
	BottomRight
)

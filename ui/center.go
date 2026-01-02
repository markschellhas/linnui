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
		// use the constrained size, this is the actual available space
		// gtx.Constraints.Max gives us the maximum space we can use
		availableSize := gtx.Constraints.Constrain(gtx.Constraints.Max)

		// measure the child with relaxed constraints (let it be its natural size)
		childGtx := gtx
		childGtx.Constraints.Min = image.Point{}
		childGtx.Constraints.Max = availableSize

		macro := op.Record(gtx.Ops)
		childDims := layout.Dimensions{}
		if child != nil {
			childDims = child(childGtx, th)
		}
		call := macro.Stop()

		// Calculate centered position
		offsetX := (availableSize.X - childDims.Size.X) / 2
		offsetY := (availableSize.Y - childDims.Size.Y) / 2

		if offsetX < 0 {
			offsetX = 0
		}
		if offsetY < 0 {
			offsetY = 0
		}

		// Apply offset and draw
		stack := op.Offset(image.Pt(offsetX, offsetY)).Push(gtx.Ops)
		call.Add(gtx.Ops)
		stack.Pop()

		return layout.Dimensions{
			Size: availableSize,
		}
	}
}

// Align positions its child according to the specified alignment
// Usage: Align(TopLeft, child)
func Align(alignment Alignment, child Widget) Widget {
	return func(gtx layout.Context, th *Theme) layout.Dimensions {
		// Take up all available space
		availableSize := gtx.Constraints.Max

		// Measure the child with relaxed constraints
		childGtx := gtx
		childGtx.Constraints.Min = image.Point{}
		childGtx.Constraints.Max = availableSize

		macro := op.Record(gtx.Ops)
		childDims := layout.Dimensions{}
		if child != nil {
			childDims = child(childGtx, th)
		}
		call := macro.Stop()

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
		stack := op.Offset(image.Pt(offsetX, offsetY)).Push(gtx.Ops)
		call.Add(gtx.Ops)
		stack.Pop()

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

package ui

import (
	"gioui.org/layout"
)

// Center centers its child within the available space
// Usage: Center(child)
func Center(child Widget) Widget {
	return func(gtx layout.Context, th *Theme) layout.Dimensions {
		return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			if child != nil {
				return child(gtx, th)
			}
			return layout.Dimensions{}
		})
	}
}

// Align positions its child according to the specified alignment
// Usage: Align(TopLeft, child)
func Align(alignment Alignment, child Widget) Widget {
	return func(gtx layout.Context, th *Theme) layout.Dimensions {
		var dir layout.Direction
		switch alignment {
		case TopLeft:
			dir = layout.NW
		case TopCenter:
			dir = layout.N
		case TopRight:
			dir = layout.NE
		case CenterLeft:
			dir = layout.W
		case CenterCenter:
			dir = layout.Center
		case CenterRight:
			dir = layout.E
		case BottomLeft:
			dir = layout.SW
		case BottomCenter:
			dir = layout.S
		case BottomRight:
			dir = layout.SE
		default:
			dir = layout.Center
		}

		return dir.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			if child != nil {
				return child(gtx, th)
			}
			return layout.Dimensions{}
		})
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

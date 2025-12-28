package ui

import (
	"gioui.org/layout"
)

// FlexWidget wraps a Widget with flex information
type FlexWidget struct {
	Widget Widget
	Flex   float32
}

// Spacer creates a flexible space that expands to fill available space
// Usage: Spacer() or Spacer(2) for flex weight of 2
func Spacer(flex ...float32) FlexWidget {
	weight := float32(1)
	if len(flex) > 0 {
		weight = flex[0]
	}

	return FlexWidget{
		Widget: func(gtx layout.Context, th *Theme) layout.Dimensions {
			// Just return the available space
			return layout.Dimensions{Size: gtx.Constraints.Min}
		},
		Flex: weight,
	}
}

// Expanded wraps a widget to make it expand to fill available space
// Usage: Expanded(child) or Expanded(child, 2) for flex weight of 2
func Expanded(child Widget, flex ...float32) FlexWidget {
	weight := float32(1)
	if len(flex) > 0 {
		weight = flex[0]
	}

	return FlexWidget{
		Widget: child,
		Flex:   weight,
	}
}

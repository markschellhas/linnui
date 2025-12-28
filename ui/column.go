package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
)

// ColumnOption configures the Column
type ColumnOption func(*columnModel)

// Spacing sets the space between children
func Spacing(dp float32) ColumnOption {
	return func(c *columnModel) { c.spacing = unit.Dp(dp) }
}

// MainAxis sets the main axis alignment for Column
func MainAxis(align MainAxisAlignment) ColumnOption {
	return func(c *columnModel) { c.mainAlign = align }
}

// columnModel holds configuration (internal)
type columnModel struct {
	spacing   unit.Dp
	mainAlign MainAxisAlignment
	children  []any // Can be Widget or FlexWidget
}

// Column creates a vertical layout
// Children can be Widget or FlexWidget (from Spacer/Expanded)
func Column(children []any, opts ...ColumnOption) Widget {
	c := &columnModel{
		spacing:   unit.Dp(8),
		mainAlign: MainAxisStart,
		children:  children,
	}
	for _, opt := range opts {
		opt(c)
	}

	return func(gtx layout.Context, th *Theme) layout.Dimensions {
		flexChildren := make([]layout.FlexChild, 0, len(c.children)*2)
		for i, child := range c.children {
			if i > 0 && c.spacing > 0 {
				flexChildren = append(flexChildren, layout.Rigid(layout.Spacer{Height: c.spacing}.Layout))
			}

			switch w := child.(type) {
			case FlexWidget:
				// Flexible child (Spacer or Expanded)
				widget := w.Widget
				flexChildren = append(flexChildren, layout.Flexed(w.Flex, func(gtx layout.Context) layout.Dimensions {
					return widget(gtx, th)
				}))
			case Widget:
				// Regular rigid child
				widget := w
				flexChildren = append(flexChildren, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return widget(gtx, th)
				}))
			}
		}
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx, flexChildren...)
	}
}

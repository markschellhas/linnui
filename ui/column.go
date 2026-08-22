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

// CrossAxis sets the cross axis alignment for Column
func CrossAxis(align CrossAxisAlignment) ColumnOption {
	return func(c *columnModel) { c.crossAlign = align }
}

// columnModel holds configuration (internal)
type columnModel struct {
	spacing    unit.Dp
	mainAlign  MainAxisAlignment
	crossAlign CrossAxisAlignment
	children   []any // Can be Widget or FlexWidget
}

// Column creates a vertical layout
// Children can be Widget or FlexWidget (from Spacer/Expanded)
func Column(children []any, opts ...ColumnOption) Widget {
	c := &columnModel{
		spacing:    unit.Dp(8),
		mainAlign:  MainAxisStart,
		crossAlign: CrossAxisStart,
		children:   children,
	}
	for _, opt := range opts {
		opt(c)
	}

	return func(gtx layout.Context, th *Theme) layout.Dimensions {
		flexChildren := make([]layout.FlexChild, 0, len(c.children))
		for _, child := range c.children {
			switch w := child.(type) {
			case FlexWidget:
				// Flexible child (Spacer or Expanded)
				widget := alignWidget(w.Widget, layout.Vertical, c.crossAlign)
				flexChildren = append(flexChildren, layout.Flexed(w.Flex, func(gtx layout.Context) layout.Dimensions {
					return widget(gtx, th)
				}))
			case Widget:
				// Regular rigid child
				widget := alignWidget(w, layout.Vertical, c.crossAlign)
				flexChildren = append(flexChildren, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return widget(gtx, th)
				}))
			}
		}
		return layout.Flex{
			Axis:      layout.Vertical,
			Spacing:   mainAxisSpacing(c.mainAlign),
			Alignment: crossAxisAlignment(c.crossAlign),
			Gap:       gtx.Dp(c.spacing),
		}.Layout(gtx, flexChildren...)
	}
}

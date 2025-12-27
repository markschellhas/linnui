package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
)

// RowOption configures the Row
type RowOption func(*rowModel)

// RowSpacing sets the space between children in a Row
func RowSpacing(dp float32) RowOption {
	return func(r *rowModel) { r.spacing = unit.Dp(dp) }
}

// RowMainAxis sets the main axis alignment for Row
func RowMainAxis(align MainAxisAlignment) RowOption {
	return func(r *rowModel) { r.mainAlign = align }
}

// rowModel holds configuration (internal)
type rowModel struct {
	spacing   unit.Dp
	mainAlign MainAxisAlignment
	children  []Widget
}

// Row creates a horizontal layout
func Row(children []Widget, opts ...RowOption) Widget {
	r := &rowModel{
		spacing:   unit.Dp(8),
		mainAlign: MainAxisStart,
		children:  children,
	}
	for _, opt := range opts {
		opt(r)
	}

	return func(gtx layout.Context, th *Theme) layout.Dimensions {
		flexChildren := make([]layout.FlexChild, 0, len(r.children)*2)
		for i, child := range r.children {
			if i > 0 {
				flexChildren = append(flexChildren, layout.Rigid(layout.Spacer{Width: r.spacing}.Layout))
			}
			w := child
			flexChildren = append(flexChildren, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return w(gtx, th)
			}))
		}
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, flexChildren...)
	}
}

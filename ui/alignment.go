package ui

import "gioui.org/layout"

// Gio uses one million pixels as the main-axis constraint for list children.
const unboundedLayoutConstraint = 1_000_000

func mainAxisSpacing(align MainAxisAlignment) layout.Spacing {
	switch align {
	case MainAxisCenter:
		return layout.SpaceSides
	case MainAxisEnd:
		return layout.SpaceStart
	case MainAxisSpaceBetween:
		return layout.SpaceBetween
	case MainAxisSpaceAround:
		return layout.SpaceAround
	case MainAxisSpaceEvenly:
		return layout.SpaceEvenly
	default:
		return layout.SpaceEnd
	}
}

func crossAxisAlignment(align CrossAxisAlignment) layout.Alignment {
	switch align {
	case CrossAxisCenter:
		return layout.Middle
	case CrossAxisEnd:
		return layout.End
	case CrossAxisBaseline:
		return layout.Baseline
	default:
		return layout.Start
	}
}

func alignWidget(widget Widget, axis layout.Axis, align CrossAxisAlignment) Widget {
	if align != CrossAxisStretch {
		return widget
	}

	return func(gtx layout.Context, th *Theme) layout.Dimensions {
		if axis == layout.Horizontal {
			if gtx.Constraints.Max.Y < unboundedLayoutConstraint {
				gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
			}
		} else {
			if gtx.Constraints.Max.X < unboundedLayoutConstraint {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
			}
		}
		return widget(gtx, th)
	}
}

package ui

import "gioui.org/layout"

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
			gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
		} else {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
		}
		return widget(gtx, th)
	}
}

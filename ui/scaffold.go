package ui

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

// ScaffoldOption configures the Scaffold
type ScaffoldOption func(*scaffoldModel)

// TopBar sets the app bar for the Scaffold.
func TopBar(bar Widget) ScaffoldOption {
	return func(s *scaffoldModel) { s.appBar = bar }
}

// Body sets the main content for the Scaffold
func Body(body Widget) ScaffoldOption {
	return func(s *scaffoldModel) { s.body = body }
}

// FAB sets the floating action button for the Scaffold
func FAB(fab Widget) ScaffoldOption {
	return func(s *scaffoldModel) { s.fab = fab }
}

// FABAlignment controls horizontal floating-action-button placement.
type FABAlignment uint8

const (
	FABEnd FABAlignment = iota
	FABCenter
	FABStart
)

// FABLocation sets horizontal floating-action-button placement.
func FABLocation(alignment FABAlignment) ScaffoldOption {
	return func(s *scaffoldModel) { s.fabAlignment = alignment }
}

// BodyAvoidsFAB reserves bottom space for a configured FAB.
func BodyAvoidsFAB(avoid bool) ScaffoldOption {
	return func(s *scaffoldModel) { s.bodyAvoidsFAB = avoid }
}

// ScaffoldOverlays adds transient or modal layers above the app content.
func ScaffoldOverlays(overlays ...Overlay) ScaffoldOption {
	return func(s *scaffoldModel) { s.overlays = append(s.overlays, overlays...) }
}

// scaffoldModel holds the configuration (internal)
type scaffoldModel struct {
	appBar        Widget
	body          Widget
	fab           Widget
	fabAlignment  FABAlignment
	bodyAvoidsFAB bool
	overlays      []Overlay
}

// Scaffold creates a top-level app layout with optional AppBar, Body, and FAB
func Scaffold(opts ...ScaffoldOption) Widget {
	s := &scaffoldModel{}
	for _, opt := range opts {
		opt(s)
	}

	return func(gtx layout.Context, th *Theme) layout.Dimensions {
		modal := false
		for _, overlay := range s.overlays {
			if overlay.active() && overlay.modal {
				modal = true
				break
			}
		}
		contentContext := gtx
		if modal {
			contentContext = gtx.Disabled()
		}

		children := []layout.StackChild{
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				size := gtx.Constraints.Min
				paint.FillShape(gtx.Ops, th.Palette.Background, clip.Rect(image.Rectangle{Max: size}).Op())
				return layout.Dimensions{Size: size}
			}),
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(contentContext,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if s.appBar != nil {
							return s.appBar(gtx, th)
						}
						return layout.Dimensions{}
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						if s.body != nil {
							bottom := unit.Dp(0)
							if s.fab != nil && s.bodyAvoidsFAB {
								bottom = unit.Dp(80)
							}
							return layout.Inset{Bottom: bottom}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return s.body(gtx, th)
							})
						}
						return layout.Dimensions{}
					}),
				)
			}),
		}

		if s.fab != nil {
			direction := layout.SE
			switch s.fabAlignment {
			case FABCenter:
				direction = layout.S
			case FABStart:
				direction = layout.SW
			}
			children = append(children, layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				return direction.Layout(contentContext, func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{
						Left: unit.Dp(16), Right: unit.Dp(16), Bottom: unit.Dp(16),
					}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return s.fab(gtx, th)
					})
				})
			}))
		}

		for _, overlay := range s.overlays {
			overlay := overlay
			if !overlay.active() {
				continue
			}
			children = append(children, layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				return overlay.widget(gtx, th)
			}))
		}
		return layout.Stack{}.Layout(gtx, children...)
	}
}

// TitleBar creates a simple title bar widget.
// Deprecated: use AppBar.
func TitleBar(title string) Widget {
	return AppBar(title)
}

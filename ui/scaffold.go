package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// ScaffoldOption configures the Scaffold
type ScaffoldOption func(*scaffoldModel)

// AppBar sets the app bar for the Scaffold
func AppBar(bar Widget) ScaffoldOption {
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

// scaffoldModel holds the configuration (internal)
type scaffoldModel struct {
	appBar Widget
	body   Widget
	fab    Widget
}

// Scaffold creates a top-level app layout with optional AppBar, Body, and FAB
func Scaffold(opts ...ScaffoldOption) Widget {
	s := &scaffoldModel{}
	for _, opt := range opts {
		opt(s)
	}

	return func(gtx layout.Context, th *Theme) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if s.appBar != nil {
					return s.appBar(gtx, th)
				}
				return layout.Dimensions{}
			}),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				if s.body != nil {
					return layout.Inset{Left: 0, Right: 0, Top: 0, Bottom: unit.Dp(80)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return s.body(gtx, th)
					})
				}
				return layout.Dimensions{}
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if s.fab != nil {
					return layout.SE.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Right: unit.Dp(16), Bottom: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return s.fab(gtx, th)
						})
					})
				}
				return layout.Dimensions{}
			}),
		)
	}
}

// TitleBar creates a simple title bar widget
func TitleBar(title string) Widget {
	return func(gtx layout.Context, th *Theme) layout.Dimensions {
		return material.List(th.Theme, new(widget.List)).Layout(gtx, 1, func(gtx layout.Context, _ int) layout.Dimensions {
			return layout.UniformInset(unit.Dp(16)).Layout(gtx, material.H6(th.Theme, title).Layout)
		})
	}
}

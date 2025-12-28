package ui

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

// ContainerOption configures the Container
type ContainerOption func(*containerModel)

// Background sets the background color
func Background(c color.NRGBA) ContainerOption {
	return func(m *containerModel) { m.background = c; m.hasBackground = true }
}

// BorderRadius sets the corner radius in dp
func BorderRadius(dp float32) ContainerOption {
	return func(m *containerModel) { m.borderRadius = dp }
}

// Border sets the border style
func Border(b BorderStyle) ContainerOption {
	return func(m *containerModel) { m.border = b; m.hasBorder = true }
}

// Shadow sets the elevation/shadow depth
func Shadow(elevation float32) ContainerOption {
	return func(m *containerModel) { m.shadow = elevation }
}

// BorderStyle defines border properties
type BorderStyle struct {
	Width float32
	Color color.NRGBA
}

// BorderAll creates a uniform border on all sides
func BorderAll(width float32, c color.NRGBA) BorderStyle {
	return BorderStyle{Width: width, Color: c}
}

// containerModel holds Container configuration (internal)
type containerModel struct {
	background    color.NRGBA
	hasBackground bool
	borderRadius  float32
	border        BorderStyle
	hasBorder     bool
	shadow        float32
}

// Container creates a decorated box that can hold a child
// Usage: Container(child, Background(White), BorderRadius(12), Shadow(4))
func Container(opts ...any) Widget {
	m := &containerModel{}
	var child Widget

	for _, opt := range opts {
		switch v := opt.(type) {
		case ContainerOption:
			v(m)
		case Widget:
			child = v
		}
	}

	return func(gtx layout.Context, th *Theme) layout.Dimensions {
		return layout.Stack{}.Layout(gtx,
			// Background and border layer
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				size := gtx.Constraints.Min
				radius := gtx.Dp(unit.Dp(m.borderRadius))

				// Draw shadow (simplified - just a darker offset rect)
				if m.shadow > 0 {
					shadowOffset := gtx.Dp(unit.Dp(m.shadow / 2))
					shadowRect := image.Rect(shadowOffset, shadowOffset, size.X+shadowOffset, size.Y+shadowOffset)
					shadowClip := clip.UniformRRect(shadowRect, radius).Push(gtx.Ops)
					paint.Fill(gtx.Ops, color.NRGBA{A: uint8(m.shadow * 8)})
					shadowClip.Pop()
				}

				// Draw background
				if m.hasBackground {
					rect := image.Rect(0, 0, size.X, size.Y)
					bgClip := clip.UniformRRect(rect, radius).Push(gtx.Ops)
					paint.Fill(gtx.Ops, m.background)
					bgClip.Pop()
				}

				// Draw border
				if m.hasBorder && m.border.Width > 0 {
					borderWidth := gtx.Dp(unit.Dp(m.border.Width))
					rect := image.Rect(0, 0, size.X, size.Y)

					// Outer clip
					outerClip := clip.UniformRRect(rect, radius).Push(gtx.Ops)
					paint.Fill(gtx.Ops, m.border.Color)
					outerClip.Pop()

					// Inner clip (punch out the inside)
					innerRect := image.Rect(borderWidth, borderWidth, size.X-borderWidth, size.Y-borderWidth)
					innerRadius := radius - borderWidth
					if innerRadius < 0 {
						innerRadius = 0
					}
					innerClip := clip.UniformRRect(innerRect, innerRadius).Push(gtx.Ops)
					if m.hasBackground {
						paint.Fill(gtx.Ops, m.background)
					} else {
						paint.Fill(gtx.Ops, color.NRGBA{A: 0})
					}
					innerClip.Pop()
				}

				return layout.Dimensions{Size: size}
			}),
			// Child layer
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				if child != nil {
					return child(gtx, th)
				}
				return layout.Dimensions{}
			}),
		)
	}
}

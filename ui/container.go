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
	return func(m *containerModel) {
		m.background = c
		m.backgroundRole = containerBackgroundCustom
	}
}

// SurfaceBackground uses the theme's surface color.
func SurfaceBackground() ContainerOption {
	return func(m *containerModel) { m.backgroundRole = containerBackgroundSurface }
}

// SurfaceVariantBackground uses the theme's surface variant color.
func SurfaceVariantBackground() ContainerOption {
	return func(m *containerModel) { m.backgroundRole = containerBackgroundSurfaceVariant }
}

// PrimaryContainerBackground uses the theme's primary container color.
func PrimaryContainerBackground() ContainerOption {
	return func(m *containerModel) { m.backgroundRole = containerBackgroundPrimary }
}

// BorderRadius sets the corner radius in dp
func BorderRadius(dp float32) ContainerOption {
	return func(m *containerModel) { m.borderRadius = dp }
}

// Border sets the border style
func Border(b BorderStyle) ContainerOption {
	return func(m *containerModel) { m.border = b; m.hasBorder = true }
}

// OutlineBorder creates a border using the theme's outline color.
func OutlineBorder(width float32) ContainerOption {
	return func(m *containerModel) {
		m.border.Width = width
		m.hasBorder = true
		m.themeOutline = true
	}
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

type containerBackground uint8

const (
	containerBackgroundNone containerBackground = iota
	containerBackgroundCustom
	containerBackgroundSurface
	containerBackgroundSurfaceVariant
	containerBackgroundPrimary
)

// containerModel holds Container configuration (internal)
type containerModel struct {
	background     color.NRGBA
	backgroundRole containerBackground
	borderRadius   float32
	border         BorderStyle
	hasBorder      bool
	themeOutline   bool
	shadow         float32
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
		background, hasBackground := containerBackgroundColor(m, th)
		borderColor := m.border.Color
		if m.themeOutline {
			borderColor = th.Palette.Outline
		}

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
					shadow := th.Palette.Shadow
					shadow.A = uint8(min(m.shadow*8, 96))
					paint.Fill(gtx.Ops, shadow)
					shadowClip.Pop()
				}

				// Draw background
				if hasBackground {
					rect := image.Rect(0, 0, size.X, size.Y)
					bgClip := clip.UniformRRect(rect, radius).Push(gtx.Ops)
					paint.Fill(gtx.Ops, background)
					bgClip.Pop()
				}

				// Draw border
				if m.hasBorder && m.border.Width > 0 {
					borderWidth := gtx.Dp(unit.Dp(m.border.Width))
					rect := image.Rect(0, 0, size.X, size.Y)
					outline := clip.Stroke{
						Path:  clip.UniformRRect(rect, radius).Path(gtx.Ops),
						Width: float32(borderWidth),
					}.Op().Push(gtx.Ops)
					paint.Fill(gtx.Ops, borderColor)
					outline.Pop()
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

func containerBackgroundColor(m *containerModel, th *Theme) (color.NRGBA, bool) {
	switch m.backgroundRole {
	case containerBackgroundCustom:
		return m.background, true
	case containerBackgroundSurface:
		return th.Palette.Surface, true
	case containerBackgroundSurfaceVariant:
		return th.Palette.SurfaceVariant, true
	case containerBackgroundPrimary:
		return th.Palette.PrimaryContainer, true
	default:
		return color.NRGBA{}, false
	}
}

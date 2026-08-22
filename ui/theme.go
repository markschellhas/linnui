package ui

import (
	"image/color"

	"gioui.org/widget/material"
)

// Common colors for easy use
var (
	White       = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	Black       = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	Transparent = color.NRGBA{A: 0}

	// Grays
	Gray50  = color.NRGBA{R: 250, G: 250, B: 250, A: 255}
	Gray100 = color.NRGBA{R: 245, G: 245, B: 245, A: 255}
	Gray200 = color.NRGBA{R: 238, G: 238, B: 238, A: 255}
	Gray300 = color.NRGBA{R: 224, G: 224, B: 224, A: 255}
	Gray400 = color.NRGBA{R: 189, G: 189, B: 189, A: 255}
	Gray500 = color.NRGBA{R: 158, G: 158, B: 158, A: 255}
	Gray600 = color.NRGBA{R: 117, G: 117, B: 117, A: 255}
	Gray700 = color.NRGBA{R: 97, G: 97, B: 97, A: 255}
	Gray800 = color.NRGBA{R: 66, G: 66, B: 66, A: 255}
	Gray900 = color.NRGBA{R: 33, G: 33, B: 33, A: 255}

	// Primary colors
	Red    = color.NRGBA{R: 244, G: 67, B: 54, A: 255}
	Pink   = color.NRGBA{R: 233, G: 30, B: 99, A: 255}
	Purple = color.NRGBA{R: 156, G: 39, B: 176, A: 255}
	Indigo = color.NRGBA{R: 63, G: 81, B: 181, A: 255}
	Blue   = color.NRGBA{R: 33, G: 150, B: 243, A: 255}
	Cyan   = color.NRGBA{R: 0, G: 188, B: 212, A: 255}
	Teal   = color.NRGBA{R: 0, G: 150, B: 136, A: 255}
	Green  = color.NRGBA{R: 76, G: 175, B: 80, A: 255}
	Yellow = color.NRGBA{R: 255, G: 235, B: 59, A: 255}
	Orange = color.NRGBA{R: 255, G: 152, B: 0, A: 255}
)

// Palette defines semantic colors used by LinnUI components.
type Palette struct {
	Primary            color.NRGBA
	OnPrimary          color.NRGBA
	PrimaryContainer   color.NRGBA
	OnPrimaryContainer color.NRGBA
	Secondary          color.NRGBA
	OnSecondary        color.NRGBA
	Background         color.NRGBA
	OnBackground       color.NRGBA
	Surface            color.NRGBA
	OnSurface          color.NRGBA
	SurfaceVariant     color.NRGBA
	OnSurfaceVariant   color.NRGBA
	Error              color.NRGBA
	OnError            color.NRGBA
	Outline            color.NRGBA
	Shadow             color.NRGBA
	Disabled           color.NRGBA
	OnDisabled         color.NRGBA
}

// Theme holds styling information
type Theme struct {
	*material.Theme
	Palette Palette
}

// ThemeFromPalette creates a Theme and keeps Gio's Material palette in sync
// with LinnUI's semantic colors.
func ThemeFromPalette(palette Palette) Theme {
	materialTheme := material.NewTheme()
	materialTheme.Palette.Bg = palette.Background
	materialTheme.Palette.Fg = palette.OnBackground
	materialTheme.Palette.ContrastBg = palette.Primary
	materialTheme.Palette.ContrastFg = palette.OnPrimary
	return Theme{Theme: materialTheme, Palette: palette}
}

// Light is the default light color scheme.
var Light = ThemeFromPalette(Palette{
	Primary:            rgb(0x6750A4),
	OnPrimary:          rgb(0xFFFFFF),
	PrimaryContainer:   rgb(0xEADDFF),
	OnPrimaryContainer: rgb(0x21005D),
	Secondary:          rgb(0x625B71),
	OnSecondary:        rgb(0xFFFFFF),
	Background:         rgb(0xFFFBFE),
	OnBackground:       rgb(0x1C1B1F),
	Surface:            rgb(0xFFFBFE),
	OnSurface:          rgb(0x1C1B1F),
	SurfaceVariant:     rgb(0xE7E0EC),
	OnSurfaceVariant:   rgb(0x49454F),
	Error:              rgb(0xB3261E),
	OnError:            rgb(0xFFFFFF),
	Outline:            rgb(0x79747E),
	Shadow:             rgb(0x000000),
	Disabled:           rgba(0x1C1B1F, 0x1F),
	OnDisabled:         rgba(0x1C1B1F, 0x61),
})

// Dark is the default dark color scheme.
var Dark = ThemeFromPalette(Palette{
	Primary:            rgb(0xD0BCFF),
	OnPrimary:          rgb(0x381E72),
	PrimaryContainer:   rgb(0x4F378B),
	OnPrimaryContainer: rgb(0xEADDFF),
	Secondary:          rgb(0xCCC2DC),
	OnSecondary:        rgb(0x332D41),
	Background:         rgb(0x1C1B1F),
	OnBackground:       rgb(0xE6E1E5),
	Surface:            rgb(0x1C1B1F),
	OnSurface:          rgb(0xE6E1E5),
	SurfaceVariant:     rgb(0x49454F),
	OnSurfaceVariant:   rgb(0xCAC4D0),
	Error:              rgb(0xF2B8B5),
	OnError:            rgb(0x601410),
	Outline:            rgb(0x938F99),
	Shadow:             rgb(0x000000),
	Disabled:           rgba(0xE6E1E5, 0x1F),
	OnDisabled:         rgba(0xE6E1E5, 0x61),
})

func rgb(hex uint32) color.NRGBA {
	return rgba(hex, 0xFF)
}

func rgba(hex uint32, alpha uint8) color.NRGBA {
	return color.NRGBA{
		R: uint8(hex >> 16),
		G: uint8(hex >> 8),
		B: uint8(hex),
		A: alpha,
	}
}

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

	// Neutral grays (shadcn / Tailwind Neutral)
	Gray50  = color.NRGBA{R: 250, G: 250, B: 250, A: 255}
	Gray100 = color.NRGBA{R: 245, G: 245, B: 245, A: 255}
	Gray200 = color.NRGBA{R: 229, G: 229, B: 229, A: 255}
	Gray300 = color.NRGBA{R: 212, G: 212, B: 212, A: 255}
	Gray400 = color.NRGBA{R: 163, G: 163, B: 163, A: 255}
	Gray500 = color.NRGBA{R: 115, G: 115, B: 115, A: 255}
	Gray600 = color.NRGBA{R: 82, G: 82, B: 82, A: 255}
	Gray700 = color.NRGBA{R: 64, G: 64, B: 64, A: 255}
	Gray800 = color.NRGBA{R: 38, G: 38, B: 38, A: 255}
	Gray900 = color.NRGBA{R: 23, G: 23, B: 23, A: 255}

	// Accent colors (Tailwind 500)
	Red    = color.NRGBA{R: 239, G: 68, B: 68, A: 255}
	Pink   = color.NRGBA{R: 236, G: 72, B: 153, A: 255}
	Purple = color.NRGBA{R: 168, G: 85, B: 247, A: 255}
	Indigo = color.NRGBA{R: 99, G: 102, B: 241, A: 255}
	Blue   = color.NRGBA{R: 59, G: 130, B: 246, A: 255}
	Cyan   = color.NRGBA{R: 6, G: 182, B: 212, A: 255}
	Teal   = color.NRGBA{R: 20, G: 184, B: 166, A: 255}
	Green  = color.NRGBA{R: 34, G: 197, B: 94, A: 255}
	Yellow = color.NRGBA{R: 234, G: 179, B: 8, A: 255}
	Orange = color.NRGBA{R: 249, G: 115, B: 22, A: 255}
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

// Light is the default light color scheme, matching shadcn Neutral.
// https://ui.shadcn.com/docs/theming
var Light = ThemeFromPalette(Palette{
	Primary:            rgb(0x171717), // oklch(0.205 0 0)
	OnPrimary:          rgb(0xFAFAFA), // oklch(0.985 0 0)
	PrimaryContainer:   rgb(0xF5F5F5), // oklch(0.97 0 0) accent / secondary
	OnPrimaryContainer: rgb(0x171717),
	Secondary:          rgb(0xF5F5F5),
	OnSecondary:        rgb(0x171717),
	Background:         rgb(0xFFFFFF), // oklch(1 0 0)
	OnBackground:       rgb(0x0A0A0A), // oklch(0.145 0 0)
	Surface:            rgb(0xFFFFFF), // card
	OnSurface:          rgb(0x0A0A0A),
	SurfaceVariant:     rgb(0xF5F5F5), // muted
	OnSurfaceVariant:   rgb(0x737373), // oklch(0.556 0 0) muted-foreground
	Error:              rgb(0xE7000B), // oklch(0.577 0.245 27.325)
	OnError:            rgb(0xFAFAFA),
	Outline:            rgb(0xE5E5E5), // oklch(0.922 0 0) border
	Shadow:             rgb(0x000000),
	Disabled:           rgba(0x0A0A0A, 0x1F),
	OnDisabled:         rgba(0x0A0A0A, 0x61),
})

// Dark is the default dark color scheme, matching shadcn Neutral.
var Dark = ThemeFromPalette(Palette{
	Primary:            rgb(0xE5E5E5), // oklch(0.922 0 0)
	OnPrimary:          rgb(0x171717),
	PrimaryContainer:   rgb(0x262626), // oklch(0.269 0 0) accent / secondary
	OnPrimaryContainer: rgb(0xFAFAFA),
	Secondary:          rgb(0x262626),
	OnSecondary:        rgb(0xFAFAFA),
	Background:         rgb(0x0A0A0A),
	OnBackground:       rgb(0xFAFAFA),
	Surface:            rgb(0x171717), // card
	OnSurface:          rgb(0xFAFAFA),
	SurfaceVariant:     rgb(0x262626), // muted
	OnSurfaceVariant:   rgb(0xA1A1A1), // oklch(0.708 0 0)
	Error:              rgb(0xFF6467), // oklch(0.704 0.191 22.216)
	OnError:            rgb(0xFAFAFA),
	Outline:            rgb(0x262626), // opaque stand-in for oklch(1 0 0 / 10%)
	Shadow:             rgb(0x000000),
	Disabled:           rgba(0xFAFAFA, 0x1F),
	OnDisabled:         rgba(0xFAFAFA, 0x61),
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

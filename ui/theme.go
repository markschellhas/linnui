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

// Palette defines the color scheme
type Palette struct {
	Primary        color.NRGBA
	OnPrimary      color.NRGBA
	SurfaceVariant color.NRGBA
	Outline        color.NRGBA
}

// Theme holds styling information
type Theme struct {
	*material.Theme
	Palette Palette
}

// Light theme with modern colors
var Light = Theme{
	Theme: material.NewTheme(),
	Palette: Palette{
		Primary:        color.NRGBA{R: 99, G: 91, B: 255, A: 255}, // Indigo
		OnPrimary:      color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		SurfaceVariant: color.NRGBA{R: 240, G: 240, B: 255, A: 255},
		Outline:        color.NRGBA{R: 150, G: 150, B: 150, A: 255},
	},
}

// Dark theme
var Dark = Theme{
	Theme: material.NewTheme(),
	Palette: Palette{
		Primary:        color.NRGBA{R: 187, G: 134, B: 252, A: 255}, // Purple
		OnPrimary:      color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		SurfaceVariant: color.NRGBA{R: 30, G: 30, B: 46, A: 255},
		Outline:        color.NRGBA{R: 100, G: 100, B: 100, A: 255},
	},
}

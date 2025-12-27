package ui

import (
	"image/color"

	"gioui.org/widget/material"
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

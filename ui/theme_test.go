package ui

import (
	"image/color"
	"testing"
)

func TestThemeFromPaletteSyncsMaterialColors(t *testing.T) {
	palette := Palette{
		Primary:      color.NRGBA{R: 1, A: 255},
		OnPrimary:    color.NRGBA{G: 2, A: 255},
		Background:   color.NRGBA{B: 3, A: 255},
		OnBackground: color.NRGBA{R: 4, G: 4, A: 255},
	}

	theme := ThemeFromPalette(palette)

	if got := theme.Theme.Palette.ContrastBg; got != palette.Primary {
		t.Errorf("material contrast background = %v, want %v", got, palette.Primary)
	}
	if got := theme.Theme.Palette.ContrastFg; got != palette.OnPrimary {
		t.Errorf("material contrast foreground = %v, want %v", got, palette.OnPrimary)
	}
	if got := theme.Theme.Palette.Bg; got != palette.Background {
		t.Errorf("material background = %v, want %v", got, palette.Background)
	}
	if got := theme.Theme.Palette.Fg; got != palette.OnBackground {
		t.Errorf("material foreground = %v, want %v", got, palette.OnBackground)
	}
}

func TestDefaultThemesHaveSemanticColors(t *testing.T) {
	tests := []struct {
		name    string
		palette Palette
	}{
		{"light", Light.Palette},
		{"dark", Dark.Palette},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			colors := map[string]color.NRGBA{
				"primary":              tt.palette.Primary,
				"on primary":           tt.palette.OnPrimary,
				"primary container":    tt.palette.PrimaryContainer,
				"on primary container": tt.palette.OnPrimaryContainer,
				"secondary":            tt.palette.Secondary,
				"on secondary":         tt.palette.OnSecondary,
				"background":           tt.palette.Background,
				"on background":        tt.palette.OnBackground,
				"surface":              tt.palette.Surface,
				"on surface":           tt.palette.OnSurface,
				"surface variant":      tt.palette.SurfaceVariant,
				"on surface variant":   tt.palette.OnSurfaceVariant,
				"error":                tt.palette.Error,
				"on error":             tt.palette.OnError,
				"outline":              tt.palette.Outline,
				"shadow":               tt.palette.Shadow,
				"disabled":             tt.palette.Disabled,
				"on disabled":          tt.palette.OnDisabled,
			}
			for name, value := range colors {
				if value.A == 0 {
					t.Errorf("%s must not be fully transparent", name)
				}
			}
		})
	}
}

func TestContainerSemanticBackgrounds(t *testing.T) {
	tests := []struct {
		name string
		role containerBackground
		want color.NRGBA
	}{
		{"surface", containerBackgroundSurface, Light.Palette.Surface},
		{"surface variant", containerBackgroundSurfaceVariant, Light.Palette.SurfaceVariant},
		{"primary container", containerBackgroundPrimary, Light.Palette.PrimaryContainer},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := containerBackgroundColor(&containerModel{backgroundRole: tt.role}, &Light)
			if !ok {
				t.Fatal("semantic background was not enabled")
			}
			if got != tt.want {
				t.Fatalf("background = %v, want %v", got, tt.want)
			}
		})
	}
}

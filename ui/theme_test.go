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

func TestDefaultThemesMatchShadcnNeutral(t *testing.T) {
	tests := []struct {
		name    string
		palette Palette
		want    Palette
	}{
		{
			name:    "light",
			palette: Light.Palette,
			want: Palette{
				Primary:            rgb(0x171717),
				OnPrimary:          rgb(0xFAFAFA),
				PrimaryContainer:   rgb(0xF5F5F5),
				OnPrimaryContainer: rgb(0x171717),
				Secondary:          rgb(0xF5F5F5),
				OnSecondary:        rgb(0x171717),
				Background:         rgb(0xFFFFFF),
				OnBackground:       rgb(0x0A0A0A),
				Surface:            rgb(0xFFFFFF),
				OnSurface:          rgb(0x0A0A0A),
				SurfaceVariant:     rgb(0xF5F5F5),
				OnSurfaceVariant:   rgb(0x737373),
				Error:              rgb(0xE7000B),
				OnError:            rgb(0xFAFAFA),
				Outline:            rgb(0xE5E5E5),
				Shadow:             rgb(0x000000),
				Disabled:           rgba(0x0A0A0A, 0x1F),
				OnDisabled:         rgba(0x0A0A0A, 0x61),
			},
		},
		{
			name:    "dark",
			palette: Dark.Palette,
			want: Palette{
				Primary:            rgb(0xE5E5E5),
				OnPrimary:          rgb(0x171717),
				PrimaryContainer:   rgb(0x262626),
				OnPrimaryContainer: rgb(0xFAFAFA),
				Secondary:          rgb(0x262626),
				OnSecondary:        rgb(0xFAFAFA),
				Background:         rgb(0x0A0A0A),
				OnBackground:       rgb(0xFAFAFA),
				Surface:            rgb(0x171717),
				OnSurface:          rgb(0xFAFAFA),
				SurfaceVariant:     rgb(0x262626),
				OnSurfaceVariant:   rgb(0xA1A1A1),
				Error:              rgb(0xFF6467),
				OnError:            rgb(0xFAFAFA),
				Outline:            rgb(0x262626),
				Shadow:             rgb(0x000000),
				Disabled:           rgba(0xFAFAFA, 0x1F),
				OnDisabled:         rgba(0xFAFAFA, 0x61),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.palette != tt.want {
				t.Errorf("palette = %+v, want %+v", tt.palette, tt.want)
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

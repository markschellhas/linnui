package ui

import (
	"math"
	"testing"
)

func TestSliderValue(t *testing.T) {
	tests := []struct {
		name            string
		value, min, max float32
		step, want      float32
	}{
		{"clamps low", -1, 0, 10, 0, 0},
		{"clamps high", 11, 0, 10, 0, 10},
		{"continuous", 4.25, 0, 10, 0, 4.25},
		{"rounds to step", 4.25, 0, 10, 0.5, 4.5},
		{"step from nonzero minimum", 7, 5, 15, 4, 9},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sliderValue(tt.value, tt.min, tt.max, tt.step); got != tt.want {
				t.Fatalf("sliderValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliderRejectsInvalidRange(t *testing.T) {
	tests := []struct {
		name     string
		min, max float32
	}{
		{"equal", 1, 1},
		{"reversed", 2, 1},
		{"nan", float32(math.NaN()), 1},
		{"infinite", 0, float32(math.Inf(1))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("Slider did not reject invalid range")
				}
			}()
			Slider("Value", NewState(float32(0)), tt.min, tt.max)
		})
	}
}

func TestControlsRequireState(t *testing.T) {
	tests := []struct {
		name string
		call func()
	}{
		{"checkbox", func() { Checkbox("Check", nil) }},
		{"switch", func() { Switch("Switch", nil) }},
		{"radio", func() { Radio(nil, nil) }},
		{"slider", func() { Slider("Value", nil, 0, 1) }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("control accepted nil state")
				}
			}()
			tt.call()
		})
	}
}

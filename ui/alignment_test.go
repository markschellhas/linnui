package ui

import (
	"image"
	"testing"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
)

func TestMainAxisSpacing(t *testing.T) {
	tests := []struct {
		name string
		in   MainAxisAlignment
		want layout.Spacing
	}{
		{"start", MainAxisStart, layout.SpaceEnd},
		{"center", MainAxisCenter, layout.SpaceSides},
		{"end", MainAxisEnd, layout.SpaceStart},
		{"space between", MainAxisSpaceBetween, layout.SpaceBetween},
		{"space around", MainAxisSpaceAround, layout.SpaceAround},
		{"space evenly", MainAxisSpaceEvenly, layout.SpaceEvenly},
		{"invalid defaults to start", MainAxisAlignment(99), layout.SpaceEnd},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mainAxisSpacing(tt.in); got != tt.want {
				t.Fatalf("mainAxisSpacing(%v) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestCrossAxisAlignment(t *testing.T) {
	tests := []struct {
		name string
		in   CrossAxisAlignment
		want layout.Alignment
	}{
		{"start", CrossAxisStart, layout.Start},
		{"center", CrossAxisCenter, layout.Middle},
		{"end", CrossAxisEnd, layout.End},
		{"stretch positioned at start", CrossAxisStretch, layout.Start},
		{"baseline", CrossAxisBaseline, layout.Baseline},
		{"invalid defaults to start", CrossAxisAlignment(99), layout.Start},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := crossAxisAlignment(tt.in); got != tt.want {
				t.Fatalf("crossAxisAlignment(%v) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestCrossAxisStretch(t *testing.T) {
	tests := []struct {
		name     string
		widget   Widget
		layout   Widget
		wantSize image.Point
	}{
		{
			name: "row stretches child height",
			widget: func(gtx layout.Context, _ *Theme) layout.Dimensions {
				if got, want := gtx.Constraints.Min, image.Pt(0, 40); got != want {
					t.Fatalf("child minimum constraints = %v, want %v", got, want)
				}
				return layout.Dimensions{Size: gtx.Constraints.Min}
			},
			layout:   nil,
			wantSize: image.Pt(0, 40),
		},
		{
			name: "column stretches child width",
			widget: func(gtx layout.Context, _ *Theme) layout.Dimensions {
				if got, want := gtx.Constraints.Min, image.Pt(100, 0); got != want {
					t.Fatalf("child minimum constraints = %v, want %v", got, want)
				}
				return layout.Dimensions{Size: gtx.Constraints.Min}
			},
			layout:   nil,
			wantSize: image.Pt(100, 0),
		},
	}

	for i := range tests {
		tt := &tests[i]
		if tt.name == "row stretches child height" {
			tt.layout = Row([]any{tt.widget}, RowCrossAxis(CrossAxisStretch))
		} else {
			tt.layout = Column([]any{tt.widget}, CrossAxis(CrossAxisStretch))
		}

		t.Run(tt.name, func(t *testing.T) {
			gtx := layout.Context{
				Ops:         new(op.Ops),
				Constraints: layout.Constraints{Max: image.Pt(100, 40)},
				Metric:      unit.Metric{PxPerDp: 1, PxPerSp: 1},
			}
			if got := tt.layout(gtx, &Light).Size; got != tt.wantSize {
				t.Fatalf("layout size = %v, want %v", got, tt.wantSize)
			}
		})
	}
}

func TestSpacingUsesFlexGap(t *testing.T) {
	var child Widget = func(_ layout.Context, _ *Theme) layout.Dimensions {
		return layout.Dimensions{Size: image.Pt(10, 10)}
	}
	gtx := layout.Context{
		Ops:         new(op.Ops),
		Constraints: layout.Constraints{Max: image.Pt(100, 100)},
		Metric:      unit.Metric{PxPerDp: 1, PxPerSp: 1},
	}

	if got, want := Row([]any{child, child}, RowSpacing(6))(gtx, &Light).Size, image.Pt(26, 10); got != want {
		t.Fatalf("row size = %v, want %v", got, want)
	}
}

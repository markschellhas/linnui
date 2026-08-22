package ui

import (
	"image/color"

	"gioui.org/layout"
)

type appBarModel struct {
	subtitle      string
	leading       Widget
	actions       []Widget
	background    color.NRGBA
	foreground    color.NRGBA
	hasBackground bool
	hasForeground bool
	elevation     float32
}

// AppBarOption configures an AppBar.
type AppBarOption func(*appBarModel)

// AppBarSubtitle adds supporting text below the title.
func AppBarSubtitle(subtitle string) AppBarOption {
	return func(model *appBarModel) { model.subtitle = subtitle }
}

// AppBarLeading places a widget before the title.
func AppBarLeading(widget Widget) AppBarOption {
	return func(model *appBarModel) { model.leading = widget }
}

// AppBarActions places widgets after the title.
func AppBarActions(actions ...Widget) AppBarOption {
	return func(model *appBarModel) { model.actions = append(model.actions, actions...) }
}

// AppBarBackground overrides the semantic surface background.
func AppBarBackground(background color.NRGBA) AppBarOption {
	return func(model *appBarModel) {
		model.background = background
		model.hasBackground = true
	}
}

// AppBarForeground overrides title and subtitle colors.
func AppBarForeground(foreground color.NRGBA) AppBarOption {
	return func(model *appBarModel) {
		model.foreground = foreground
		model.hasForeground = true
	}
}

// AppBarElevation sets the shadow depth in dp.
func AppBarElevation(elevation float32) AppBarOption {
	return func(model *appBarModel) { model.elevation = max(0, elevation) }
}

// AppBar creates a responsive top application bar.
func AppBar(title string, opts ...AppBarOption) Widget {
	model := appBarModel{elevation: 1}
	for _, opt := range opts {
		opt(&model)
	}

	return func(gtx layout.Context, th *Theme) layout.Dimensions {
		foreground := th.Palette.OnSurface
		if model.hasForeground {
			foreground = model.foreground
		}
		titleChildren := []any{Text(title, Style(H6), TextColor(foreground))}
		if model.subtitle != "" {
			titleChildren = append(titleChildren, Text(model.subtitle, Style(Caption), TextColor(foreground)))
		}
		rowChildren := make([]any, 0, len(model.actions)+3)
		if model.leading != nil {
			rowChildren = append(rowChildren, model.leading, SizedBox(Width(12)))
		}
		rowChildren = append(rowChildren, Expanded(Column(titleChildren, Spacing(2))))
		for _, action := range model.actions {
			rowChildren = append(rowChildren, action)
		}

		content := Padding(
			InsetsSymmetric(16, 8),
			Row(rowChildren, RowSpacing(8), RowCrossAxis(CrossAxisCenter)),
		)
		options := []any{content, BorderRadius(0), Shadow(model.elevation), SurfaceBackground()}
		if model.hasBackground {
			options = append(options, Background(model.background))
		}
		minHeight := gtx.Dp(64)
		if gtx.Constraints.Min.Y < minHeight {
			gtx.Constraints.Min.Y = min(minHeight, gtx.Constraints.Max.Y)
		}
		return Container(options...)(gtx, th)
	}
}

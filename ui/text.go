package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

// TextStyle defines preset text styles
type TextStyle int

const (
	BodyText TextStyle = iota
	H1
	H2
	H3
	H4
	H5
	H6
	Caption
	Overline
)

// TextOption configures the Text widget
type TextOption func(*textModel)

// Size sets a custom font size in sp
func Size(sp float32) TextOption {
	return func(t *textModel) { t.size = unit.Sp(sp) }
}

// Style sets a preset text style (H1, H2, Body, etc.)
func Style(s TextStyle) TextOption {
	return func(t *textModel) { t.style = s }
}

// textModel holds text configuration (internal)
type textModel struct {
	content string
	style   TextStyle
	size    unit.Sp // custom size overrides style
}

// Text creates a text display widget
func Text(content string, opts ...TextOption) Widget {
	t := &textModel{
		content: content,
		style:   BodyText,
		size:    0, // 0 means use style default
	}
	for _, opt := range opts {
		opt(t)
	}

	return func(gtx layout.Context, th *Theme) layout.Dimensions {
		var label material.LabelStyle

		// Apply preset style
		switch t.style {
		case H1:
			label = material.H1(th.Theme, t.content)
		case H2:
			label = material.H2(th.Theme, t.content)
		case H3:
			label = material.H3(th.Theme, t.content)
		case H4:
			label = material.H4(th.Theme, t.content)
		case H5:
			label = material.H5(th.Theme, t.content)
		case H6:
			label = material.H6(th.Theme, t.content)
		case Caption:
			label = material.Caption(th.Theme, t.content)
		case Overline:
			label = material.Overline(th.Theme, t.content)
		default:
			label = material.Body1(th.Theme, t.content)
		}

		// Override with custom size if specified
		if t.size > 0 {
			label.TextSize = t.size
		}

		return label.Layout(gtx)
	}
}

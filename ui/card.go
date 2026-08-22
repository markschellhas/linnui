package ui

import (
	"gioui.org/io/semantic"
	"gioui.org/layout"
)

// CardVariant controls a Card's visual emphasis.
type CardVariant uint8

const (
	CardElevated CardVariant = iota
	CardFilled
	CardOutlined
)

type cardModel struct {
	id          string
	variant     CardVariant
	padding     Insets
	radius      float32
	elevation   float32
	onClick     func()
	description string
	disabled    bool
}

// CardOption configures a Card.
type CardOption func(*cardModel)

// CardStyle sets the visual variant.
func CardStyle(variant CardVariant) CardOption {
	return func(model *cardModel) { model.variant = variant }
}

// CardPadding sets the content padding.
func CardPadding(insets Insets) CardOption {
	return func(model *cardModel) { model.padding = insets }
}

// CardRadius sets the corner radius in dp.
func CardRadius(radius float32) CardOption {
	return func(model *cardModel) { model.radius = max(0, radius) }
}

// CardElevation sets elevated-card shadow depth.
func CardElevation(elevation float32) CardOption {
	return func(model *cardModel) { model.elevation = max(0, elevation) }
}

// CardOnClick makes the entire Card interactive.
func CardOnClick(callback func()) CardOption {
	return func(model *cardModel) { model.onClick = callback }
}

// CardID sets the persistent interaction-state ID.
func CardID(id string) CardOption {
	return func(model *cardModel) { model.id = id }
}

// CardDescription adds an accessible description to an interactive Card.
func CardDescription(description string) CardOption {
	return func(model *cardModel) { model.description = description }
}

// CardDisabled controls whether an interactive Card receives input.
func CardDisabled(disabled bool) CardOption {
	return func(model *cardModel) { model.disabled = disabled }
}

// Card creates a polished content surface.
func Card(child Widget, opts ...CardOption) Widget {
	return legacyTree.Card(child, opts...)
}

// Card creates a Card whose interaction state belongs to the Tree.
func (t *Tree) Card(child Widget, opts ...CardOption) Widget {
	model := cardModel{
		id:        "card_default",
		variant:   CardElevated,
		padding:   InsetsAll(16),
		radius:    12,
		elevation: 4,
	}
	for _, opt := range opts {
		opt(&model)
	}

	content := Padding(model.padding, child)
	var clickable interface {
		Clicked(layout.Context) bool
		Layout(layout.Context, layout.Widget) layout.Dimensions
	}
	if model.onClick != nil {
		clickable = t.clickable(cardWidget, model.id)
	}

	return func(gtx layout.Context, th *Theme) layout.Dimensions {
		decorated := func(gtx layout.Context, th *Theme) layout.Dimensions {
			options := []any{content, BorderRadius(model.radius)}
			switch model.variant {
			case CardFilled:
				options = append(options, SurfaceVariantBackground())
			case CardOutlined:
				options = append(options, SurfaceBackground(), OutlineBorder(1))
			default:
				options = append(options, SurfaceBackground(), Shadow(model.elevation))
			}
			return Container(options...)(gtx, th)
		}
		if clickable == nil {
			return decorated(gtx, th)
		}
		if model.disabled {
			gtx = gtx.Disabled()
		}
		for clickable.Clicked(gtx) {
			model.onClick()
		}
		return clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			semantic.Button.Add(gtx.Ops)
			if model.description != "" {
				semantic.DescriptionOp(model.description).Add(gtx.Ops)
			}
			return decorated(gtx, th)
		})
	}
}

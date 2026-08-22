package ui

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// TextFieldOption configures the TextField
type TextFieldOption func(*textFieldModel)

// Hint sets the placeholder text
func Hint(text string) TextFieldOption {
	return func(t *textFieldModel) { t.hint = text }
}

// OnChange sets the callback when text changes
func OnChange(fn func(string)) TextFieldOption {
	return func(t *textFieldModel) { t.onChange = fn }
}

// TextFieldID sets a unique ID for the text field (for state persistence)
func TextFieldID(id string) TextFieldOption {
	return func(t *textFieldModel) { t.id = id }
}

// MultiLine allows multiple lines of text
func MultiLine() TextFieldOption {
	return func(t *textFieldModel) { t.multiLine = true }
}

// BindText binds the field bidirectionally to reactive state.
func BindText(value *State[string]) TextFieldOption {
	return func(t *textFieldModel) { t.value = value }
}

// TextFieldDisabled controls whether the field can receive input.
func TextFieldDisabled(disabled bool) TextFieldOption {
	return func(t *textFieldModel) { t.disabled = disabled }
}

// textFieldModel holds text field configuration (internal)
type textFieldModel struct {
	id        string
	hint      string
	onChange  func(string)
	multiLine bool
	value     *State[string]
	disabled  bool
}

// TextField creates a text input widget
// Usage: TextField(Hint("Enter name"), OnChange(func(s string) { ... }))
func TextField(opts ...TextFieldOption) Widget {
	return legacyTree.TextField(opts...)
}

// TextField creates a text input whose editor state belongs to the Tree.
func (tree *Tree) TextField(opts ...TextFieldOption) Widget {
	model := &textFieldModel{
		id:   "textfield_default",
		hint: "",
	}
	for _, opt := range opts {
		opt(model)
	}

	// Get persistent editor using the ID
	fieldState := tree.textField(model.id)
	editor := &fieldState.editor
	editor.SingleLine = !model.multiLine
	onChange := model.onChange

	return func(gtx layout.Context, th *Theme) layout.Dimensions {
		if model.disabled {
			gtx = gtx.Disabled()
		}

		// Check for text changes
		changed := false
		for {
			event, ok := editor.Update(gtx)
			if !ok {
				break
			}
			if _, ok := event.(widget.ChangeEvent); ok {
				changed = true
			}
		}

		if changed {
			value := editor.Text()
			if model.value != nil {
				model.value.Set(value)
			}
			if onChange != nil {
				onChange(value)
			}
		} else {
			fieldState.mu.Lock()
			pending := fieldState.pending
			fieldState.pending = nil
			fieldState.mu.Unlock()

			target := ""
			syncValue := false
			if pending != nil {
				target = *pending
				syncValue = true
				if model.value != nil {
					model.value.Set(target)
				}
			} else if model.value != nil {
				target = model.value.Get()
				syncValue = target != editor.Text()
			}
			if syncValue && target != editor.Text() {
				editor.SetText(target)
				for {
					if _, ok := editor.Update(gtx); !ok {
						break
					}
				}
			}
		}

		fieldState.mu.Lock()
		fieldState.snapshot = editor.Text()
		fieldState.mu.Unlock()

		// Style the text field
		ed := material.Editor(th.Theme, editor, model.hint)
		ed.TextSize = unit.Sp(16)
		ed.Color = th.Palette.OnSurface
		ed.HintColor = th.Palette.OnSurfaceVariant
		ed.SelectionColor = th.Palette.PrimaryContainer

		return layout.Stack{}.Layout(gtx,
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				size := gtx.Constraints.Min
				rect := image.Rectangle{Max: size}
				radius := gtx.Dp(unit.Dp(12))

				background := clip.UniformRRect(rect, radius).Push(gtx.Ops)
				paint.Fill(gtx.Ops, th.Palette.SurfaceVariant)
				background.Pop()

				outline := clip.Stroke{
					Path:  clip.UniformRRect(rect, radius).Path(gtx.Ops),
					Width: float32(gtx.Dp(unit.Dp(1))),
				}.Op().Push(gtx.Ops)
				paint.Fill(gtx.Ops, th.Palette.Outline)
				outline.Pop()
				return layout.Dimensions{Size: size}
			}),
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(unit.Dp(12)).Layout(gtx, ed.Layout)
			}),
		)
	}
}

// TextFieldValue gets the current text value for a TextField by ID
func TextFieldValue(id string) string {
	value, ok := legacyTree.TextFieldValue(id)
	if ok {
		return value
	}
	return ""
}

// SetTextFieldValue sets the text value for a TextField by ID
func SetTextFieldValue(id string, text string) {
	legacyTree.SetTextFieldValue(id, text)
}

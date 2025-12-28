package ui

import (
	"sync"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// editorRegistry stores editor state by ID
var (
	editorRegistry = make(map[string]*widget.Editor)
	editorMu       sync.Mutex
)

// getEditor returns a persistent editor for the given ID
func getEditor(id string) *widget.Editor {
	editorMu.Lock()
	defer editorMu.Unlock()

	if e, ok := editorRegistry[id]; ok {
		return e
	}
	e := new(widget.Editor)
	e.SingleLine = true
	editorRegistry[id] = e
	return e
}

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

// textFieldModel holds text field configuration (internal)
type textFieldModel struct {
	id        string
	hint      string
	onChange  func(string)
	multiLine bool
}

// TextField creates a text input widget
// Usage: TextField(Hint("Enter name"), OnChange(func(s string) { ... }))
func TextField(opts ...TextFieldOption) Widget {
	t := &textFieldModel{
		id:   "textfield_default",
		hint: "",
	}
	for _, opt := range opts {
		opt(t)
	}

	// Get persistent editor using the ID
	editor := getEditor(t.id)
	editor.SingleLine = !t.multiLine
	onChange := t.onChange // Capture the handler

	return func(gtx layout.Context, th *Theme) layout.Dimensions {
		// Check for text changes
		for {
			event, ok := editor.Update(gtx)
			if !ok {
				break
			}
			if _, ok := event.(widget.ChangeEvent); ok {
				if onChange != nil {
					onChange(editor.Text())
				}
			}
		}

		// Style the text field
		ed := material.Editor(th.Theme, editor, t.hint)
		ed.TextSize = unit.Sp(16)

		// Wrap in a bordered container
		return layout.Background{}.Layout(gtx,
			func(gtx layout.Context) layout.Dimensions {
				return layout.Dimensions{Size: gtx.Constraints.Min}
			},
			func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(unit.Dp(12)).Layout(gtx, ed.Layout)
			},
		)
	}
}

// TextFieldValue gets the current text value for a TextField by ID
func TextFieldValue(id string) string {
	editorMu.Lock()
	defer editorMu.Unlock()

	if e, ok := editorRegistry[id]; ok {
		return e.Text()
	}
	return ""
}

// SetTextFieldValue sets the text value for a TextField by ID
func SetTextFieldValue(id string, text string) {
	editorMu.Lock()
	defer editorMu.Unlock()

	if e, ok := editorRegistry[id]; ok {
		e.SetText(text)
	}
}

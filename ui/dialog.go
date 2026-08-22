package ui

import (
	"image"
	"strconv"
	"sync"

	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
)

// DialogState controls the visibility of a modal Dialog.
type DialogState struct {
	mu          sync.RWMutex
	visible     bool
	invalidator Invalidator
	scrim       widget.Clickable
	tree        *Tree
}

// NewDialogState creates hidden dialog state.
func NewDialogState() *DialogState {
	return &DialogState{tree: NewTree(nil)}
}

// Bind attaches dialog updates to a render host.
func (state *DialogState) Bind(invalidator Invalidator) *DialogState {
	state.mu.Lock()
	state.invalidator = invalidator
	state.mu.Unlock()
	return state
}

// Show displays the dialog.
func (state *DialogState) Show() {
	state.setVisible(true)
}

// Dismiss hides the dialog.
func (state *DialogState) Dismiss() {
	state.setVisible(false)
}

// Visible reports whether the dialog is displayed.
func (state *DialogState) Visible() bool {
	state.mu.RLock()
	defer state.mu.RUnlock()
	return state.visible
}

func (state *DialogState) setVisible(visible bool) {
	state.mu.Lock()
	changed := state.visible != visible
	state.visible = visible
	invalidator := state.invalidator
	state.mu.Unlock()
	if changed && invalidator != nil {
		invalidator.Invalidate()
	}
}

// DialogAction describes one action button.
type DialogAction struct {
	Label    string
	Variant  ButtonVariant
	OnClick  func()
	KeepOpen bool
}

type dialogModel struct {
	title           string
	content         Widget
	actions         []DialogAction
	description     string
	dismissOnScrim  bool
	dismissOnEscape bool
}

// DialogOption configures a Dialog.
type DialogOption func(*dialogModel)

// DialogTitle sets the dialog heading.
func DialogTitle(title string) DialogOption {
	return func(model *dialogModel) { model.title = title }
}

// DialogContent sets the main dialog content.
func DialogContent(content Widget) DialogOption {
	return func(model *dialogModel) { model.content = content }
}

// DialogActions sets action buttons.
func DialogActions(actions ...DialogAction) DialogOption {
	return func(model *dialogModel) { model.actions = append(model.actions, actions...) }
}

// DialogDescription provides additional accessible context.
func DialogDescription(description string) DialogOption {
	return func(model *dialogModel) { model.description = description }
}

// DialogDismissOnScrim controls dismissal when the backdrop is clicked.
func DialogDismissOnScrim(enabled bool) DialogOption {
	return func(model *dialogModel) { model.dismissOnScrim = enabled }
}

// DialogDismissOnEscape controls dismissal from the Escape key.
func DialogDismissOnEscape(enabled bool) DialogOption {
	return func(model *dialogModel) { model.dismissOnEscape = enabled }
}

// Dialog creates a modal Scaffold overlay.
func Dialog(state *DialogState, opts ...DialogOption) Overlay {
	if state == nil {
		panic("ui: Dialog requires non-nil DialogState")
	}
	model := dialogModel{dismissOnScrim: true, dismissOnEscape: true}
	for _, opt := range opts {
		opt(&model)
	}

	widget := func(gtx layout.Context, th *Theme) layout.Dimensions {
		gtx.Constraints.Min = gtx.Constraints.Max
		if model.dismissOnEscape {
			for {
				event, ok := gtx.Event(key.Filter{Name: key.NameEscape})
				if !ok {
					break
				}
				if event, ok := event.(key.Event); ok && event.State == key.Press {
					state.Dismiss()
				}
			}
		}
		for state.scrim.Clicked(gtx) {
			if model.dismissOnScrim {
				state.Dismiss()
			}
		}

		return layout.Stack{}.Layout(gtx,
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				return state.scrim.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					size := gtx.Constraints.Min
					scrim := th.Palette.Shadow
					scrim.A = 150
					paint.FillShape(gtx.Ops, scrim, clip.Rect(image.Rectangle{Max: size}).Op())
					return layout.Dimensions{Size: size}
				})
			}),
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					margin := gtx.Dp(unit.Dp(32))
					width := min(gtx.Dp(unit.Dp(480)), max(0, gtx.Constraints.Max.X-margin))
					gtx.Constraints.Min.X = width
					gtx.Constraints.Max.X = width

					children := make([]any, 0, 3)
					if model.title != "" {
						children = append(children, Text(model.title, Style(H5)))
					}
					if model.content != nil {
						children = append(children, model.content)
					}
					if len(model.actions) > 0 {
						actions := make([]any, 0, len(model.actions))
						for index, action := range model.actions {
							action := action
							actions = append(actions, state.tree.Button(
								action.Label,
								ButtonID(dialogActionID(index, action.Label)),
								Variant(action.Variant),
								OnClick(func() {
									if action.OnClick != nil {
										action.OnClick()
									}
									if !action.KeepOpen {
										state.Dismiss()
									}
								}),
							))
						}
						children = append(children, Row(actions, RowSpacing(8), RowMainAxis(MainAxisEnd)))
					}
					return state.tree.Card(
						Column(children, Spacing(16)),
						CardID("dialog"),
						CardStyle(CardElevated),
						CardPadding(InsetsAll(24)),
						CardDescription(model.description),
					)(gtx, th)
				})
			}),
		)
	}
	return CustomOverlay(widget, OverlayVisible(state.Visible), OverlayModal())
}

func dialogActionID(index int, label string) string {
	return "dialog_action_" + strconv.Itoa(index) + "_" + label
}

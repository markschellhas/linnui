package ui

import (
	"image"
	"strconv"
	"sync"

	"gioui.org/io/event"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
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
	content     widget.Clickable
	modalTag    struct{}
	cardBounds  image.Rectangle
	scrimID     pointer.ID
	scrimDown   bool
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
		for state.content.Clicked(gtx) {
			// Consume clicks on non-interactive dialog content so they cannot
			// reach controls behind the modal.
		}
		for {
			input, ok := gtx.Event(pointer.Filter{
				Target: &state.modalTag,
				Kinds:  pointer.Press | pointer.Release | pointer.Move | pointer.Drag | pointer.Scroll | pointer.Cancel,
			})
			if !ok {
				break
			}
			event, ok := input.(pointer.Event)
			if !ok {
				continue
			}
			switch event.Kind {
			case pointer.Press:
				position := image.Pt(int(event.Position.X), int(event.Position.Y))
				state.mu.Lock()
				outside := !state.cardBounds.Contains(position)
				if outside {
					state.scrimDown = true
					state.scrimID = event.PointerID
				}
				state.mu.Unlock()
				if outside {
					gtx.Execute(pointer.GrabCmd{Tag: &state.modalTag, ID: event.PointerID})
				}
			case pointer.Release:
				state.mu.Lock()
				dismiss := state.scrimDown && state.scrimID == event.PointerID
				state.scrimDown = false
				state.mu.Unlock()
				if dismiss && model.dismissOnScrim {
					state.Dismiss()
				}
			case pointer.Cancel:
				state.mu.Lock()
				state.scrimDown = false
				state.mu.Unlock()
			}
		}

		fullSize := gtx.Constraints.Max
		modalBounds := clip.Rect(image.Rectangle{Max: fullSize}).Push(gtx.Ops)
		event.Op(gtx.Ops, &state.modalTag)
		margin := gtx.Dp(unit.Dp(32))
		width := min(gtx.Dp(unit.Dp(480)), max(0, fullSize.X-margin))
		cardContext := gtx
		cardContext.Constraints = layout.Constraints{
			Min: image.Pt(width, 0),
			Max: image.Pt(width, max(0, fullSize.Y-margin)),
		}
		recording := op.Record(gtx.Ops)
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
		cardDimensions := state.tree.Card(
			Column(children, Spacing(16)),
			CardID("dialog"),
			CardStyle(CardElevated),
			CardPadding(InsetsAll(24)),
			CardDescription(model.description),
		)(cardContext, th)
		cardCall := recording.Stop()
		cardOffset := image.Pt(
			(fullSize.X-cardDimensions.Size.X)/2,
			(fullSize.Y-cardDimensions.Size.Y)/2,
		)

		scrim := th.Palette.Shadow
		scrim.A = 150
		paint.FillShape(gtx.Ops, scrim, clip.Rect(image.Rectangle{Max: fullSize}).Op())

		cardBounds := image.Rectangle{Min: cardOffset, Max: cardOffset.Add(cardDimensions.Size)}
		state.mu.Lock()
		state.cardBounds = cardBounds
		state.mu.Unlock()

		offset := op.Offset(cardOffset).Push(gtx.Ops)
		contentContext := gtx
		contentContext.Constraints = layout.Exact(cardDimensions.Size)
		state.content.Layout(contentContext, func(gtx layout.Context) layout.Dimensions {
			return layout.Dimensions{Size: gtx.Constraints.Min}
		})
		cardCall.Add(gtx.Ops)
		offset.Pop()
		modalBounds.Pop()
		return layout.Dimensions{Size: fullSize}
	}
	return CustomOverlay(widget, OverlayVisible(state.Visible), OverlayModal())
}

func dialogActionID(index int, label string) string {
	return "dialog_action_" + strconv.Itoa(index) + "_" + label
}

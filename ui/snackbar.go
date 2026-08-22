package ui

import (
	"sync"
	"time"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
)

const defaultSnackbarDuration = 4 * time.Second

// SnackbarState controls transient Snackbar messages.
type SnackbarState struct {
	mu          sync.RWMutex
	message     string
	deadline    time.Time
	visible     bool
	invalidator Invalidator
	tree        *Tree
}

// NewSnackbarState creates hidden Snackbar state.
func NewSnackbarState() *SnackbarState {
	return &SnackbarState{tree: NewTree(nil)}
}

// Bind attaches Snackbar updates to a render host.
func (state *SnackbarState) Bind(invalidator Invalidator) *SnackbarState {
	state.mu.Lock()
	state.invalidator = invalidator
	state.mu.Unlock()
	return state
}

// Show displays a message for the default duration.
func (state *SnackbarState) Show(message string) {
	state.ShowFor(message, defaultSnackbarDuration)
}

// ShowFor displays a message for duration. A non-positive duration persists
// until explicitly dismissed.
func (state *SnackbarState) ShowFor(message string, duration time.Duration) {
	state.mu.Lock()
	state.message = message
	state.visible = message != ""
	state.deadline = time.Time{}
	if duration > 0 {
		state.deadline = time.Now().Add(duration)
	}
	invalidator := state.invalidator
	state.mu.Unlock()
	if invalidator != nil {
		invalidator.Invalidate()
	}
}

// Dismiss hides the current message.
func (state *SnackbarState) Dismiss() {
	state.mu.Lock()
	changed := state.visible
	state.visible = false
	invalidator := state.invalidator
	state.mu.Unlock()
	if changed && invalidator != nil {
		invalidator.Invalidate()
	}
}

// Visible reports whether a message is active.
func (state *SnackbarState) Visible() bool {
	state.mu.RLock()
	defer state.mu.RUnlock()
	return state.visible
}

func (state *SnackbarState) current() (string, time.Time, bool) {
	state.mu.RLock()
	defer state.mu.RUnlock()
	return state.message, state.deadline, state.visible
}

type snackbarModel struct {
	actionLabel string
	action      func()
}

// SnackbarOption configures a Snackbar.
type SnackbarOption func(*snackbarModel)

// SnackbarAction adds an optional action button.
func SnackbarAction(label string, callback func()) SnackbarOption {
	return func(model *snackbarModel) {
		model.actionLabel = label
		model.action = callback
	}
}

// Snackbar creates a transient Scaffold overlay.
func Snackbar(state *SnackbarState, opts ...SnackbarOption) Overlay {
	if state == nil {
		panic("ui: Snackbar requires non-nil SnackbarState")
	}
	model := snackbarModel{}
	for _, opt := range opts {
		opt(&model)
	}

	content := func(gtx layout.Context, th *Theme) layout.Dimensions {
		message, deadline, visible := state.current()
		if !visible {
			return layout.Dimensions{}
		}
		if !deadline.IsZero() {
			if !gtx.Now.Before(deadline) {
				state.Dismiss()
				return layout.Dimensions{}
			}
			gtx.Execute(op.InvalidateCmd{At: deadline})
		}

		children := []any{Expanded(Text(message, TextColor(th.Palette.OnSurfaceVariant)))}
		if model.actionLabel != "" {
			children = append(children, state.tree.Button(
				model.actionLabel,
				ButtonID("snackbar_action"),
				Variant(TextButton),
				OnClick(func() {
					if model.action != nil {
						model.action()
					}
					state.Dismiss()
				}),
			))
		}
		maxWidth := gtx.Dp(unit.Dp(600))
		if gtx.Constraints.Max.X > maxWidth {
			gtx.Constraints.Max.X = maxWidth
		}
		return Container(
			Padding(
				InsetsSymmetric(16, 10),
				Row(children, RowSpacing(12), RowCrossAxis(CrossAxisCenter)),
			),
			SurfaceVariantBackground(),
			BorderRadius(12),
			Shadow(6),
		)(gtx, th)
	}

	return CustomOverlay(
		alignedOverlay(layout.S, InsetsOnly(0, 16, 16, 16), content),
		OverlayVisible(state.Visible),
	)
}

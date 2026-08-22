package ui

import (
	"fmt"
	"math"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

type boolControlModel struct {
	id       string
	disabled bool
	onChange func(bool)
}

// CheckboxOption configures a Checkbox.
type CheckboxOption func(*boolControlModel)

// CheckboxID sets the persistent state ID.
func CheckboxID(id string) CheckboxOption {
	return func(model *boolControlModel) { model.id = id }
}

// CheckboxDisabled controls whether the checkbox can receive input.
func CheckboxDisabled(disabled bool) CheckboxOption {
	return func(model *boolControlModel) { model.disabled = disabled }
}

// OnCheckboxChange runs after a user changes the value.
func OnCheckboxChange(callback func(bool)) CheckboxOption {
	return func(model *boolControlModel) { model.onChange = callback }
}

// Checkbox creates a checkbox bound to reactive state.
func Checkbox(label string, value *State[bool], opts ...CheckboxOption) Widget {
	return legacyTree.Checkbox(label, value, opts...)
}

// Checkbox creates a checkbox whose interaction state belongs to the Tree.
func (t *Tree) Checkbox(label string, value *State[bool], opts ...CheckboxOption) Widget {
	if value == nil {
		panic("ui: Checkbox requires non-nil State[bool]")
	}
	model := boolControlModel{id: label}
	for _, opt := range opts {
		opt(&model)
	}
	control := t.boolControl(checkboxWidget, model.id)

	return func(gtx layout.Context, th *Theme) layout.Dimensions {
		control.Value = value.Get()
		before := control.Value
		if model.disabled {
			gtx = gtx.Disabled()
		}
		style := material.CheckBox(th.Theme, control, label)
		dims := style.Layout(gtx)
		if !model.disabled && control.Value != before {
			value.Set(control.Value)
			if model.onChange != nil {
				model.onChange(control.Value)
			}
		}
		return dims
	}
}

// SwitchOption configures a Switch.
type SwitchOption func(*boolControlModel)

// SwitchID sets the persistent state ID.
func SwitchID(id string) SwitchOption {
	return func(model *boolControlModel) { model.id = id }
}

// SwitchDisabled controls whether the switch can receive input.
func SwitchDisabled(disabled bool) SwitchOption {
	return func(model *boolControlModel) { model.disabled = disabled }
}

// OnSwitchChange runs after a user changes the value.
func OnSwitchChange(callback func(bool)) SwitchOption {
	return func(model *boolControlModel) { model.onChange = callback }
}

// Switch creates a labeled switch bound to reactive state.
func Switch(label string, value *State[bool], opts ...SwitchOption) Widget {
	return legacyTree.Switch(label, value, opts...)
}

// Switch creates a switch whose interaction state belongs to the Tree.
func (t *Tree) Switch(label string, value *State[bool], opts ...SwitchOption) Widget {
	if value == nil {
		panic("ui: Switch requires non-nil State[bool]")
	}
	model := boolControlModel{id: label}
	for _, opt := range opts {
		opt(&model)
	}
	control := t.boolControl(switchWidget, model.id)

	return func(gtx layout.Context, th *Theme) layout.Dimensions {
		control.Value = value.Get()
		before := control.Value
		if model.disabled {
			gtx = gtx.Disabled()
		}
		style := material.Switch(th.Theme, control, label)
		style.Color.Enabled = th.Palette.Primary
		style.Color.Disabled = th.Palette.SurfaceVariant
		style.Color.Track = th.Palette.Outline

		textColor := th.Palette.OnSurface
		if model.disabled {
			textColor = th.Palette.OnDisabled
		}
		dims := layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return Text(label, TextColor(textColor))(gtx, th)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return style.Layout(gtx)
			}),
		)
		if !model.disabled && control.Value != before {
			value.Set(control.Value)
			if model.onChange != nil {
				model.onChange(control.Value)
			}
		}
		return dims
	}
}

// RadioChoice is one value and visible label in a Radio group.
type RadioChoice struct {
	Value string
	Label string
}

type radioModel struct {
	id       string
	disabled bool
	onChange func(string)
}

// RadioOption configures a Radio group.
type RadioOption func(*radioModel)

// RadioID sets the persistent state ID.
func RadioID(id string) RadioOption {
	return func(model *radioModel) { model.id = id }
}

// RadioDisabled controls whether the group can receive input.
func RadioDisabled(disabled bool) RadioOption {
	return func(model *radioModel) { model.disabled = disabled }
}

// OnRadioChange runs after a user selects a value.
func OnRadioChange(callback func(string)) RadioOption {
	return func(model *radioModel) { model.onChange = callback }
}

// Radio creates a radio group bound to reactive string state.
func Radio(value *State[string], choices []RadioChoice, opts ...RadioOption) Widget {
	return legacyTree.Radio(value, choices, opts...)
}

// Radio creates a radio group whose interaction state belongs to the Tree.
func (t *Tree) Radio(value *State[string], choices []RadioChoice, opts ...RadioOption) Widget {
	if value == nil {
		panic("ui: Radio requires non-nil State[string]")
	}
	model := radioModel{id: "radio_default"}
	for _, opt := range opts {
		opt(&model)
	}
	group := t.enum(model.id)

	return func(gtx layout.Context, th *Theme) layout.Dimensions {
		group.Value = value.Get()
		before := group.Value
		if model.disabled {
			gtx = gtx.Disabled()
		}
		children := make([]layout.FlexChild, 0, len(choices))
		for _, choice := range choices {
			choice := choice
			children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return material.RadioButton(th.Theme, group, choice.Value, choice.Label).Layout(gtx)
			}))
		}
		dims := layout.Flex{Axis: layout.Vertical, Gap: gtx.Dp(unit.Dp(4))}.Layout(gtx, children...)
		if !model.disabled && group.Value != before {
			value.Set(group.Value)
			if model.onChange != nil {
				model.onChange(group.Value)
			}
		}
		return dims
	}
}

type sliderModel struct {
	id        string
	disabled  bool
	step      float32
	vertical  bool
	onChange  func(float32)
	valueText func(float32) string
}

// SliderOption configures a Slider.
type SliderOption func(*sliderModel)

// SliderID sets the persistent state ID.
func SliderID(id string) SliderOption {
	return func(model *sliderModel) { model.id = id }
}

// SliderDisabled controls whether the slider can receive input.
func SliderDisabled(disabled bool) SliderOption {
	return func(model *sliderModel) { model.disabled = disabled }
}

// SliderStep snaps values to a positive interval.
func SliderStep(step float32) SliderOption {
	return func(model *sliderModel) { model.step = step }
}

// SliderVertical lays out the slider vertically.
func SliderVertical() SliderOption {
	return func(model *sliderModel) { model.vertical = true }
}

// OnSliderChange runs after a user changes the value.
func OnSliderChange(callback func(float32)) SliderOption {
	return func(model *sliderModel) { model.onChange = callback }
}

// SliderValueText formats the displayed value.
func SliderValueText(format func(float32) string) SliderOption {
	return func(model *sliderModel) { model.valueText = format }
}

// Slider creates a labeled range control bound to reactive state.
func Slider(label string, value *State[float32], minValue, maxValue float32, opts ...SliderOption) Widget {
	return legacyTree.Slider(label, value, minValue, maxValue, opts...)
}

// Slider creates a range control whose interaction state belongs to the Tree.
func (t *Tree) Slider(label string, value *State[float32], minValue, maxValue float32, opts ...SliderOption) Widget {
	if value == nil {
		panic("ui: Slider requires non-nil State[float32]")
	}
	if !finite(minValue) || !finite(maxValue) || minValue >= maxValue {
		panic("ui: Slider requires finite min < max")
	}
	model := sliderModel{
		id:        label,
		valueText: func(value float32) string { return fmt.Sprintf("%.0f", value) },
	}
	for _, opt := range opts {
		opt(&model)
	}
	if model.step < 0 || !finite(model.step) {
		panic("ui: Slider step must be finite and non-negative")
	}
	control := t.float(model.id)

	return func(gtx layout.Context, th *Theme) layout.Dimensions {
		current := sliderValue(value.Get(), minValue, maxValue, model.step)
		control.Value = (current - minValue) / (maxValue - minValue)
		before := control.Value
		if model.disabled {
			gtx = gtx.Disabled()
		}

		axis := layout.Horizontal
		if model.vertical {
			axis = layout.Vertical
			gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
		} else {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
		}
		style := material.Slider(th.Theme, control)
		style.Axis = axis
		style.Color = th.Palette.Primary
		slider := Widget(func(gtx layout.Context, _ *Theme) layout.Dimensions {
			return style.Layout(gtx)
		})

		textColor := th.Palette.OnSurface
		if model.disabled {
			textColor = th.Palette.OnDisabled
		}
		var dims layout.Dimensions
		if model.vertical {
			dims = layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return Text(label, TextColor(textColor))(gtx, th)
				}),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return slider(gtx, th)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return Text(model.valueText(current), TextColor(textColor))(gtx, th)
				}),
			)
		} else {
			dims = Column([]any{
				Row([]any{
					Text(label, TextColor(textColor)),
					Spacer(),
					Text(model.valueText(current), TextColor(textColor)),
				}, RowSpacing(8)),
				slider,
			}, Spacing(4))(gtx, th)
		}

		if !model.disabled && control.Value != before {
			next := sliderValue(minValue+control.Value*(maxValue-minValue), minValue, maxValue, model.step)
			value.Set(next)
			if model.onChange != nil {
				model.onChange(next)
			}
		}
		return dims
	}
}

func sliderValue(value, minValue, maxValue, step float32) float32 {
	value = max(minValue, min(value, maxValue))
	if step > 0 {
		value = minValue + float32(math.Round(float64((value-minValue)/step)))*step
		value = max(minValue, min(value, maxValue))
	}
	return value
}

func finite(value float32) bool {
	return !math.IsNaN(float64(value)) && !math.IsInf(float64(value), 0)
}

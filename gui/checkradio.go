// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"github.com/g3n/engine/gui/assets"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/window"
)

const (
	checkON  = string(assets.CheckBox)
	checkOFF = string(assets.CheckBoxOutlineBlank)
	radioON  = string(assets.RadioButtonChecked)
	radioOFF = string(assets.RadioButtonUnchecked)
)

type CheckRadio struct {
	Panel             // Embedded panel
	Label      *Label // Text label
	icon       *Label
	styles     *CheckRadioStyles
	check      bool
	group      string
	cursorOver bool
	state      bool
	codeON     string
	codeOFF    string
}

type CheckRadioStyle struct {
	Border      BorderSizes
	Paddings    BorderSizes
	BorderColor math32.Color4
	BgColor     math32.Color4
	FgColor     math32.Color
}

type CheckRadioStyles struct {
	Normal   CheckRadioStyle
	Over     CheckRadioStyle
	Focus    CheckRadioStyle
	Disabled CheckRadioStyle
}

// NewCheckBox creates and returns a pointer to a new CheckBox widget
// with the specified text
func NewCheckBox(text string) *CheckRadio {

	return newCheckRadio(true, text)
}

// NewRadioButton creates and returns a pointer to a new RadioButton widget
// with the specified text
func NewRadioButton(text string) *CheckRadio {

	return newCheckRadio(false, text)
}

// newCheckRadio creates and returns a pointer to a new CheckRadio widget
// with the specified type and text
func newCheckRadio(check bool, text string) *CheckRadio {

	cb := new(CheckRadio)
	cb.styles = &StyleDefault.CheckRadio

	// Adapts to specified type: CheckBox or RadioButton
	cb.check = check
	cb.state = false
	if cb.check {
		cb.codeON = checkON
		cb.codeOFF = checkOFF
	} else {
		cb.codeON = radioON
		cb.codeOFF = radioOFF
	}

	// Initialize panel
	cb.Panel.Initialize(0, 0)

	// Subscribe to events
	cb.Panel.Subscribe(OnKeyDown, cb.onKey)
	cb.Panel.Subscribe(OnCursorEnter, cb.onCursor)
	cb.Panel.Subscribe(OnCursorLeave, cb.onCursor)
	cb.Panel.Subscribe(OnMouseDown, cb.onMouse)
	cb.Panel.Subscribe(OnEnable, func(evname string, ev interface{}) { cb.update() })

	// Creates label
	cb.Label = NewLabel(text)
	cb.Label.Subscribe(OnResize, func(evname string, ev interface{}) { cb.recalc() })
	cb.Panel.Add(cb.Label)

	// Creates icon label
	cb.icon = NewIconLabel(" ")
	cb.Panel.Add(cb.icon)

	cb.recalc()
	cb.update()
	return cb
}

// SetRoot overrides the IPanel.SetRoot method
func (cb *CheckRadio) SetRoot(root *Root) {

	if cb.root == root {
		return
	}
	cb.root = root
	// Subscribes once to this root panel OnRadioGroup events
	root.Subscribe(OnRadioGroup, func(name string, ev interface{}) {
		cb.onRadioGroup(ev.(*CheckRadio))
	})
}

// Value returns the current state of the checkbox
func (cb *CheckRadio) Value() bool {

	return cb.state
}

// SetValue sets the current state of the checkbox
func (cb *CheckRadio) SetValue(state bool) *CheckRadio {

	if state == cb.state {
		return cb
	}
	cb.state = state
	cb.update()
	return cb
}

// Group returns the name of the radio group
func (cb *CheckRadio) Group() string {

	return cb.group
}

// SetGroup sets the name of the radio group
func (cb *CheckRadio) SetGroup(group string) {

	cb.group = group
}

// SetStyles set the button styles overriding the default style
func (cb *CheckRadio) SetStyles(bs *CheckRadioStyles) {

	cb.styles = bs
	cb.update()
}

// toggleState toggles the current state of the checkbox/radiobutton
func (cb *CheckRadio) toggleState() {

	if cb.check {
		cb.state = !cb.state
	} else {
		if len(cb.group) == 0 {
			cb.state = !cb.state
		} else {
			if cb.state {
				return
			}
			cb.state = !cb.state
		}
	}
	cb.update()
	cb.Dispatch(OnChange, nil)
	if !cb.check && len(cb.group) > 0 {
		cb.root.Dispatch(OnRadioGroup, cb)
	}
}

// onMouse process OnMouseDown events
func (cb *CheckRadio) onMouse(evname string, ev interface{}) {

	cb.root.SetKeyFocus(cb)
	cb.root.StopPropagation(Stop3D)
	cb.toggleState()
	// Dispatch OnClick for left mouse button down
	if evname == OnMouseDown {
		mev := ev.(*window.MouseEvent)
		if mev.Button == window.MouseButtonLeft {
			cb.Dispatch(OnClick, nil)
		}
	}
	cb.root.StopPropagation(StopAll)
}

// onCursor process OnCursor* events
func (cb *CheckRadio) onCursor(evname string, ev interface{}) {

	if evname == OnCursorEnter {
		cb.cursorOver = true
	} else {
		cb.cursorOver = false
	}
	cb.update()
	cb.root.StopPropagation(StopAll)
}

// onKey receives subscribed key events
func (cb *CheckRadio) onKey(evname string, ev interface{}) {

	kev := ev.(*window.KeyEvent)
	if evname == OnKeyDown && kev.Keycode == window.KeyEnter {
		cb.toggleState()
		cb.update()
		cb.Dispatch(OnClick, nil)
		cb.root.StopPropagation(Stop3D)
		return
	}
	return
}

// onRadioGroup receives subscriber OnRadioGroup events
func (cb *CheckRadio) onRadioGroup(other *CheckRadio) {

	// If event is for this button, ignore
	if cb == other {
		return
	}
	// If other radio group is not the group of this button, ignore
	if cb.group != other.group {
		return
	}
	// Toggle this button state
	cb.SetValue(!other.Value())
}

// update updates the visual appearance of the checkbox
func (cb *CheckRadio) update() {

	if cb.state {
		cb.icon.SetText(cb.codeON)
	} else {
		cb.icon.SetText(cb.codeOFF)
	}

	if !cb.Enabled() {
		cb.applyStyle(&cb.styles.Disabled)
		return
	}
	if cb.cursorOver {
		cb.applyStyle(&cb.styles.Over)
		return
	}
	cb.applyStyle(&cb.styles.Normal)
}

// setStyle sets the specified checkradio style
func (cb *CheckRadio) applyStyle(s *CheckRadioStyle) {

	cb.Panel.SetBordersColor4(&s.BorderColor)
	cb.Panel.SetBordersFrom(&s.Border)
	cb.Panel.SetPaddingsFrom(&s.Paddings)
	cb.Panel.SetColor4(&s.BgColor)

	cb.icon.SetColor(&s.FgColor)
	cb.Label.SetColor(&s.FgColor)
}

// recalc recalculates dimensions and position from inside out
func (cb *CheckRadio) recalc() {

	// Sets icon position
	cb.icon.SetFontSize(cb.Label.FontSize() * 1.3)
	cb.icon.SetPosition(0, 0)

	// Label position
	spacing := float32(4)
	cb.Label.SetPosition(cb.icon.Width()+spacing, 0)

	// Content width
	width := cb.icon.Width() + spacing + cb.Label.Width()
	cb.SetContentSize(width, cb.Label.Height())
}

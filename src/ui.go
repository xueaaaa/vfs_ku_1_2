package main

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const (
	windowTitle      = "VFS"
	windowWidth      = 700
	windowHeight     = 450
	inputPlaceholder = "Enter command here..."
	lineBreak        = "\n"
)

type console struct {
	label  *widget.Label
	scroll *container.Scroll
	input  *widget.Entry
	view   fyne.CanvasObject
}

func newConsole() *console {
	label := widget.NewLabel("")
	label.Wrapping = fyne.TextWrapWord
	label.Selectable = true

	scroll := container.NewVScroll(label)

	input := widget.NewEntry()
	input.SetPlaceHolder(inputPlaceholder)

	return &console{
		label:  label,
		scroll: scroll,
		input:  input,
		view:   container.NewBorder(nil, input, nil, nil, scroll),
	}
}

func (c *console) Print(lines ...string) {
	var text strings.Builder
	text.WriteString(c.label.Text)
	for _, line := range lines {
		text.WriteString(line + lineBreak)
	}
	c.label.SetText(text.String())
	c.scroll.ScrollToBottom()
}

func (c *console) OnSubmit(handler func(string)) {
	c.input.OnSubmitted = func(text string) {
		handler(text)
		c.input.SetText("")
	}
}

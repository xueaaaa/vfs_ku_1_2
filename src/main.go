package main

import (
	"fmt"
	"strings"
	"vfs/src/internal/command"
	"vfs/src/internal/executor"
	"vfs/src/internal/parser"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const WindowName = "VFS"

func main() {
	a := app.New()
	w := a.NewWindow(WindowName)
	w.Resize(fyne.NewSize(700, 450))

	output := widget.NewMultiLineEntry()
	output.Disable()
	appendOutput := func(lines ...string) {
		text := output.Text
		for _, l := range lines {
			text += l + "\n"
		}
		output.SetText(text)
	}

	command.NewCommand(
		"ls",
		[]string{"arg1", "arg2"},
		func(args []string) (any, error) {
			appendOutput(fmt.Sprintf("command = ls, args = %s", args))
			return nil, nil
		},
	)

	command.NewCommand(
		"cd",
		[]string{"arg1", "arg2"},
		func(args []string) (any, error) {
			appendOutput(fmt.Sprintf("command = cd, args = %s", args))
			return nil, nil
		},
	)

	command.NewCommand(
		"exit",
		[]string{},
		func(strings []string) (any, error) {
			a.Quit()
			return nil, nil
		},
	)

	input := widget.NewEntry()
	input.SetPlaceHolder("Enter command here...")

	runCommand := func(line string) {
		line = strings.TrimSpace(line)
		appendOutput("> " + line)

		cmd, args, err := parser.Parse(line)
		if err != nil {
			appendOutput(err.Error())
			return
		}

		_, err = executor.Execute(cmd, args)
		if err != nil {
			appendOutput(err.Error())
			return
		}
	}

	input.OnSubmitted = func(text string) {
		runCommand(text)
		input.SetText("")
	}

	content := container.NewBorder(nil, input, nil, nil, container.NewScroll(output))
	w.SetContent(content)

	w.ShowAndRun()
}

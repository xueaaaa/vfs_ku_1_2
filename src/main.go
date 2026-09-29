package main

import (
	"strings"
	"vfs/src/internal/command"
	"vfs/src/internal/config"
	"vfs/src/internal/executor"
	"vfs/src/internal/parser"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

const (
	exitName   = "exit"
	promptMark = "> "
)

type printer interface {
	Print(lines ...string)
	Clear()
}

func main() {
	a := app.New()
	w := a.NewWindow(windowTitle)
	w.Resize(fyne.NewSize(windowWidth, windowHeight))

	c := newConsole()
	registerExit(a)
	run := newRunner(c)
	c.OnSubmit(run)

	w.SetContent(c.view)
	w.Show()
	startup(config.New(), c, run)
	a.Run()
}

func newRunner(out printer) func(string) {
	return func(line string) {
		line = strings.TrimSpace(line)
		out.Print(promptMark + line)

		cmd, args, err := parser.Parse(line)
		if err == nil {
			_, err = executor.Execute(cmd, args)
		}
		if err != nil {
			out.Print(err.Error())
		}
	}
}

func registerExit(a fyne.App) {
	command.NewCommand(exitName, []string{}, true, func([]string) (any, error) {
		a.Quit()
		return nil, nil
	})
}

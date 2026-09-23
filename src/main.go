package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"vfs/src/internal/command"
	"vfs/src/internal/config"
	"vfs/src/internal/executor"
	"vfs/src/internal/parser"
	"vfs/src/internal/session"
	"vfs/src/internal/vfs"

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
		"exit",
		[]string{},
		true,
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
	w.Show()
	loadConfig(appendOutput, runCommand)
	a.Run()
}

func loadConfig(appendOutput func(...string), runCommand func(line string)) {
	conf := config.New()
	if conf.VfsPath != "" {
		appendOutput(fmt.Sprintf("vfs path is set as %s", conf.VfsPath))

		vfs, err := vfs.Load(conf.VfsPath)
		if err != nil {
			appendOutput(err.Error())
		} else {
			s := session.NewSession(&vfs.Root)

			loadCommands(*s, appendOutput)

			appendOutput(fmt.Sprintf("vfs %s loaded.\nRoot folder name: "+
				"%s\nRoot subdirectories: %s\nRoot files: %s",
				vfs.Name,
				vfs.Root.Name,
				vfs.Root.Subdirs,
				vfs.Root.Files,
			),
			)
		}
	}
	if conf.ScriptPath != "" {
		appendOutput(fmt.Sprintf("script path is set as %s", conf.ScriptPath))

		f, err := os.Open(conf.ScriptPath)
		if err != nil {
			appendOutput(err.Error())
			return
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "//") {
				continue
			}
			runCommand(line)
		}

		if err := scanner.Err(); err != nil {
			appendOutput(err.Error())
		}
	}
}

func loadCommands(session session.Session, appendOutput func(...string)) {
	command.NewCommand(
		"ls",
		[]string{},
		false,
		func(args []string) (any, error) {
			display := ""
			for _, sd := range session.Current.Subdirs {
				display += sd.Name + " "
			}
			for _, f := range session.Current.Files {
				display += f.Name + " "
			}
			if display == "" {
				return nil, fmt.Errorf("No files or directories in %s\n", session.Current.Name)
			}
			appendOutput(display)
			return nil, nil
		},
	)

	command.NewCommand(
		"cd",
		nil,
		true,
		func(args []string) (any, error) {
			if len(args) != 1 {
				return nil, errors.New("invalid count of arguments")
			}
			if err := session.Cd(args[0]); err != nil {
				return nil, err
			}
			appendOutput(fmt.Sprintf("current directory is %s", session.PathString()))
			return nil, nil
		},
	)
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"vfs/src/internal/builtin"
	"vfs/src/internal/config"
	"vfs/src/internal/session"
	"vfs/src/internal/vfs"
)

const commentPrefix = "//"

func startup(conf config.Config, out printer, run func(string)) {
	if conf.VfsPath != "" {
		out.Print(fmt.Sprintf("vfs path is set as %s", conf.VfsPath))
		loadVfs(conf.VfsPath, out)
	}
	if conf.ScriptPath != "" {
		runScript(conf.ScriptPath, out, run)
	}
}

func loadVfs(path string, out printer) {
	loaded, err := vfs.Load(path)
	if err != nil {
		out.Print(err.Error())
		return
	}

	builtin.Register(session.NewSession(&loaded.Root), out)

	out.Print(
		fmt.Sprintf("vfs %s loaded.", loaded.Name),
		fmt.Sprintf("Root folder name: %s", loaded.Root.Name),
		fmt.Sprintf("Root subdirectories: %s", dirNames(loaded.Root)),
		fmt.Sprintf("Root files: %s", fileNames(loaded.Root)),
	)
}

func dirNames(dir vfs.Directory) string {
	names := make([]string, 0, len(dir.Subdirs))
	for _, d := range dir.Subdirs {
		names = append(names, d.Name)
	}
	return strings.Join(names, " ")
}

func fileNames(dir vfs.Directory) string {
	names := make([]string, 0, len(dir.Files))
	for _, f := range dir.Files {
		names = append(names, f.Name)
	}
	return strings.Join(names, " ")
}

func runScript(path string, out printer, run func(string)) {
	out.Print(fmt.Sprintf("script path is set as %s", path))

	f, err := os.Open(path)
	if err != nil {
		out.Print(err.Error())
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if line := strings.TrimSpace(scanner.Text()); !isSkipped(line) {
			run(line)
		}
	}
	if err := scanner.Err(); err != nil {
		out.Print(err.Error())
	}
}

func isSkipped(line string) bool {
	return line == "" || strings.HasPrefix(line, commentPrefix)
}

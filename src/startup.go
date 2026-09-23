package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"vfs/src/internal/config"
)

const commentPrefix = "//"

func startup(conf config.Config, out printer, run func(string)) {
	if conf.VfsPath != "" {
		out.Print(fmt.Sprintf("vfs path is set as %s", conf.VfsPath))
	}
	if conf.ScriptPath != "" {
		runScript(conf.ScriptPath, out, run)
	}
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

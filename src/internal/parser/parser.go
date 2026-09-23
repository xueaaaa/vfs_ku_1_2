package parser

import (
	"errors"
	"os"
	"strings"
)

const (
	commandIndex = 0
	argsStart    = 1
)

var errEmptyCommand = errors.New("given command is empty")

func Parse(input string) (string, []string, error) {
	line := strings.TrimSpace(input)
	if line == "" {
		return "", nil, errEmptyCommand
	}

	parts := strings.Fields(os.ExpandEnv(line))
	if len(parts) == 0 {
		return "", nil, errEmptyCommand
	}

	return parts[commandIndex], parts[argsStart:], nil
}

package parser

import (
	"errors"
	"os"
	"strings"
)

func Parse(input string) (string, []string, error) {
	if input == "" {
		return "", nil, errors.New("given command is empty")
	}

	line := strings.TrimSpace(input)
	expanded := os.ExpandEnv(line)

	parts := strings.Fields(expanded)
	cmd := parts[0]
	args := make([]string, 0)
	if len(parts) > 1 {
		args = parts[1:]
	}

	return cmd, args, nil
}

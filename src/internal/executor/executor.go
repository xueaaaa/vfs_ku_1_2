package executor

import (
	"fmt"
	"vfs/src/internal/command"
)

func Execute(cmd string, args []string) (any, error) {
	command, ok := command.Find(cmd)
	if !ok {
		return nil, fmt.Errorf("command %s not found", cmd)
	}

	return command.Execute(args)
}

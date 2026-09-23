package executor

import (
	"fmt"
	"vfs/src/internal/command"
)

func Execute(name string, args []string) (any, error) {
	cmd, ok := command.Find(name)
	if !ok {
		return nil, fmt.Errorf("command %s not found", name)
	}

	return cmd.Execute(args)
}

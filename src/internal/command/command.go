package command

import (
	"fmt"
	"slices"
)

type Command struct {
	Name          string
	AvailableArgs []string
	Execute       func([]string) (any, error)
}

var Commands []Command

func NewCommand(
	name string,
	availableArgs []string,
	execute func([]string) (any, error),
) {
	Commands = append(Commands, Command{
		Name:          name,
		AvailableArgs: availableArgs,
		Execute: func(args []string) (any, error) {
			for _, arg := range args {
				if !slices.Contains(availableArgs, arg) {
					return nil, fmt.Errorf("unknown argument: %s", arg)
				}
			}

			return execute(args)
		},
	})
}

func Find(name string) (Command, bool) {
	idx := slices.IndexFunc(Commands, func(c Command) bool {
		return c.Name == name
	})
	if idx == -1 {
		return Command{}, false
	}
	return Commands[idx], true
}

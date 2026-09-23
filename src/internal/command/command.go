package command

import (
	"fmt"
	"slices"
)

type Command struct {
	Name          string
	AvailableArgs []string
	WithoutArgs   bool
	Execute       func([]string) (any, error)
}

var Commands []Command

func NewCommand(
	name string,
	availableArgs []string,
	withoutArgs bool,
	execute func([]string) (any, error),
) {
	Commands = append(Commands, Command{
		Name:          name,
		AvailableArgs: availableArgs,
		WithoutArgs:   withoutArgs,
		Execute: func(args []string) (any, error) {
			for _, arg := range args {
				if !slices.Contains(availableArgs, arg) && !withoutArgs {
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

package builtin

import (
	"errors"
	"fmt"

	"vfs/src/internal/session"
)

const (
	cdName     = "cd"
	cdArgCount = 1

	msgCurrentDir = "current directory is %s"
)

func newCd(s *session.Session, out Console) func([]string) (any, error) {
	return func(args []string) (any, error) {
		if len(args) != cdArgCount {
			return nil, errors.New(msgBadArgCount)
		}
		if err := s.Cd(args[argFirst]); err != nil {
			return nil, err
		}
		out.Print(fmt.Sprintf(msgCurrentDir, s.PathString()))
		return nil, nil
	}
}

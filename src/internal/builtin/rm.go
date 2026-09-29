package builtin

import (
	"errors"
	"fmt"

	"vfs/src/internal/session"
)

const (
	rmName = "rm"

	msgRmNoName = "name not specified"
	msgDeleted  = "%s %s deleted"
)

func newRm(s *session.Session, out Console) func([]string) (any, error) {
	return func(args []string) (any, error) {
		if len(args) == 0 {
			return nil, errors.New(msgRmNoName)
		}
		var errs []error
		for _, target := range args {
			kind, err := s.Remove(target)
			if err != nil {
				errs = append(errs, err)
				continue
			}
			out.Print(fmt.Sprintf(msgDeleted, kind, target))
		}
		return nil, errors.Join(errs...)
	}
}

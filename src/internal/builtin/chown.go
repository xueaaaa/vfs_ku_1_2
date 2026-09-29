package builtin

import (
	"errors"
	"fmt"

	"vfs/src/internal/session"
)

const (
	chownName     = "chown"
	chownArgCount = 2
	chownOwnerArg = argFirst
	chownFileArg  = argSecond

	msgChownUsage   = "chown <owner> <file>"
	msgOwnerChanged = "owner of %s is now %s"
)

func newChown(s *session.Session, out Console) func([]string) (any, error) {
	return func(args []string) (any, error) {
		if len(args) != chownArgCount {
			return nil, errors.New(msgChownUsage)
		}
		owner, target := args[chownOwnerArg], args[chownFileArg]

		file, err := s.FindFile(target)
		if err != nil {
			return nil, err
		}
		file.Owner = owner
		out.Print(fmt.Sprintf(msgOwnerChanged, target, owner))
		return nil, nil
	}
}

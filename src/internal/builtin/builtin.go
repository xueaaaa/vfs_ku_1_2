package builtin

import (
	"vfs/src/internal/command"
	"vfs/src/internal/session"
)

const (
	argFirst = 0

	msgBadArgCount = "invalid count of arguments"
)

type Console interface {
	Print(lines ...string)
}

func Register(s *session.Session, out Console) {
	command.NewCommand(lsName, nil, true, newLs(s, out))
	command.NewCommand(cdName, nil, true, newCd(s, out))
}

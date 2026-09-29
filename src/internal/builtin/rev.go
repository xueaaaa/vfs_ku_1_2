package builtin

import (
	"errors"
	"slices"

	"vfs/src/internal/session"
)

const (
	revName     = "rev"
	revArgCount = 1
)

func newRev(s *session.Session, out Console) func([]string) (any, error) {
	return func(args []string) (any, error) {
		if len(args) != revArgCount {
			return nil, errors.New(msgBadArgCount)
		}
		file, err := s.FindFile(args[argFirst])
		if err != nil {
			return nil, err
		}
		result := reverse(string(file.Content))
		out.Print(result)
		return result, nil
	}
}

func reverse(text string) string {
	runes := []rune(text)
	slices.Reverse(runes)
	return string(runes)
}

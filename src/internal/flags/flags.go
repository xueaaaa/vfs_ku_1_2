package flags

import (
	"fmt"
	"strings"
)

const prefix = "-"

type Set map[rune]bool

func (s Set) Has(letter rune) bool { return s[letter] }

func Parse(args []string, allowed string) (Set, []string, error) {
	set := Set{}
	operands := make([]string, 0, len(args))

	for _, arg := range args {
		if !isFlag(arg) {
			operands = append(operands, arg)
			continue
		}
		if err := set.add(arg, allowed); err != nil {
			return nil, nil, err
		}
	}
	return set, operands, nil
}

func isFlag(arg string) bool {
	return strings.HasPrefix(arg, prefix) && arg != prefix
}

func (s Set) add(arg, allowed string) error {
	for _, letter := range strings.TrimPrefix(arg, prefix) {
		if !strings.ContainsRune(allowed, letter) {
			return fmt.Errorf("unknown flag: -%c", letter)
		}
		s[letter] = true
	}
	return nil
}

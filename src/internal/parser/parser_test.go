package parser

import (
	"slices"
	"testing"
)

func TestParse(t *testing.T) {
	cmd, args, err := Parse("  ls  -l   home ")
	if err != nil || cmd != "ls" || !slices.Equal(args, []string{"-l", "home"}) {
		t.Fatalf("got %q %v %v", cmd, args, err)
	}
}

func TestParseEmptyAfterExpansion(t *testing.T) {
	t.Setenv("VFS_EMPTY_VAR", "")
	for _, line := range []string{"", "   ", "$VFS_EMPTY_VAR"} {
		if _, _, err := Parse(line); err == nil {
			t.Errorf("Parse(%q): want error, got nil", line)
		}
	}
}

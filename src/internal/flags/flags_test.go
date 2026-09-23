package flags

import (
	"slices"
	"testing"
)

func TestParseCombined(t *testing.T) {
	cases := map[string]string{"-la": "la", "-hl": "hl", "-lah": "lah"}
	for arg, letters := range cases {
		set, _, err := Parse([]string{arg}, "lah")
		if err != nil || len(set) != len(letters) {
			t.Fatalf("%q: set=%v err=%v", arg, set, err)
		}
		for _, letter := range letters {
			if !set.Has(letter) {
				t.Errorf("%q: flag %c missing", arg, letter)
			}
		}
	}
}

func TestParseOperands(t *testing.T) {
	set, ops, err := Parse([]string{"home", "-l", "/etc", "-"}, "l")
	if err != nil || !set.Has('l') || !slices.Equal(ops, []string{"home", "/etc", "-"}) {
		t.Fatalf("set=%v ops=%v err=%v", set, ops, err)
	}
}

func TestParseUnknownFlag(t *testing.T) {
	if _, _, err := Parse([]string{"-lz"}, "l"); err == nil {
		t.Fatal("want error for -z")
	}
}

package builtin

import "testing"

func TestRevPaths(t *testing.T) {
	s, out := newEnv(t)
	_ = s.Cd("/home")
	if err := run(t, s, out, "rev", "user/notes.txt"); err != nil || out.text() != "tnetnoc tset" {
		t.Fatalf("rev: %q %v", out.text(), err)
	}
	if err := run(t, s, out, "rev", "/nope.txt"); err == nil {
		t.Error("missing file must fail")
	}
}

func TestClear(t *testing.T) {
	s, out := newEnv(t)
	out.Print("junk")
	_ = run(t, s, out, "clear")
	if out.text() != "console cleared" {
		t.Errorf("got %q", out.text())
	}
}

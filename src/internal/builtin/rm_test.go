package builtin

import (
	"strings"
	"testing"
)

func TestRmPathsAfterCd(t *testing.T) {
	s, out := newEnv(t)
	_ = run(t, s, out, "cd", "/home/user")
	if err := run(t, s, out, "rm", "documents/deep.txt"); err != nil {
		t.Fatalf("relative rm: %v", err)
	}
	if err := run(t, s, out, "rm", "/readme.txt"); err != nil {
		t.Fatalf("absolute rm: %v", err)
	}
	if err := run(t, s, out, "rm", "/home/user/documents/deep.txt"); err == nil {
		t.Error("second rm must fail: already deleted")
	}
	out.Clear()
	_ = run(t, s, out, "ls", "documents", "/")
	if strings.Contains(out.text(), "deep.txt") || strings.Contains(out.text(), "readme.txt") {
		t.Errorf("files still listed:\n%s", out.text())
	}
}

func TestRmDirKeepsCurrentValid(t *testing.T) {
	s, out := newEnv(t)
	_ = run(t, s, out, "cd", "/etc")
	if err := run(t, s, out, "rm", "/home"); err != nil {
		t.Fatal(err)
	}
	out.Clear()
	_ = run(t, s, out, "ls")
	if out.text() != "config.txt" {
		t.Errorf("current dir broken after rm: %q", out.text())
	}
}

func TestRmRefusals(t *testing.T) {
	s, out := newEnv(t)
	_ = run(t, s, out, "cd", "/home/user")
	for _, target := range []string{"/home", "..", ".", "/home/user", "/"} {
		if err := run(t, s, out, "rm", target); err == nil {
			t.Errorf("rm %s must fail", target)
		}
	}
	if s.PathString() != "/home/user" {
		t.Errorf("session moved: %s", s.PathString())
	}
}

func TestRmMultiple(t *testing.T) {
	s, out := newEnv(t)
	err := run(t, s, out, "rm", "readme.txt", "nope", "etc")
	if err == nil || !strings.Contains(err.Error(), "nope") {
		t.Errorf("want error about nope, got %v", err)
	}
	if !strings.Contains(out.text(), "file readme.txt deleted") ||
		!strings.Contains(out.text(), "directory etc deleted") {
		t.Errorf("output: %q", out.text())
	}
}

func TestChownAndLong(t *testing.T) {
	s, out := newEnv(t)
	_ = run(t, s, out, "cd", "/etc")
	if err := run(t, s, out, "chown", "user1", "/home/user/notes.txt"); err != nil {
		t.Fatal(err)
	}
	out.Clear()
	_ = run(t, s, out, "ls", "-l", "/home/user")
	if !strings.Contains(out.text(), "- user1") {
		t.Errorf("owner not shown:\n%s", out.text())
	}
	for _, args := range [][]string{{"user1"}, {"a", "b", "c"}, {"user1", "nowhere.txt"}} {
		if err := run(t, s, out, "chown", args...); err == nil {
			t.Errorf("chown %v must fail", args)
		}
	}
}

package builtin

import (
	"strings"
	"testing"

	"vfs/src/internal/command"
	"vfs/src/internal/session"
	"vfs/src/internal/vfs"
)

const nestedVfs = "../../../tests/stage3/vfs_nested.xml"

type fakeConsole struct{ lines []string }

func (f *fakeConsole) Print(lines ...string) { f.lines = append(f.lines, lines...) }
func (f *fakeConsole) Clear()                { f.lines = nil }
func (f *fakeConsole) text() string          { return strings.Join(f.lines, "\n") }

func run(t *testing.T, s *session.Session, out *fakeConsole, name string, args ...string) error {
	t.Helper()
	command.Commands = nil
	Register(s, out)
	cmd, ok := command.Find(name)
	if !ok {
		t.Fatalf("command %s is not registered", name)
	}
	_, err := cmd.Execute(args)
	return err
}

func newEnv(t *testing.T) (*session.Session, *fakeConsole) {
	t.Helper()
	loaded, err := vfs.Load(nestedVfs)
	if err != nil {
		t.Fatal(err)
	}
	return session.NewSession(&loaded.Root), &fakeConsole{}
}

func TestLsPlainAndPath(t *testing.T) {
	s, out := newEnv(t)
	if err := run(t, s, out, "ls"); err != nil || out.text() != "home etc readme.txt" {
		t.Fatalf("ls: %q %v", out.text(), err)
	}
	out.Clear()
	if err := run(t, s, out, "ls", "home/user"); err != nil || out.text() != "documents notes.txt" {
		t.Fatalf("ls home/user: %q %v", out.text(), err)
	}
	out.Clear()
	_ = s.Cd("/home/user/documents")
	if err := run(t, s, out, "ls", "/etc"); err != nil || out.text() != "config.txt" {
		t.Fatalf("ls /etc: %q %v", out.text(), err)
	}
}

func TestLsLongHuman(t *testing.T) {
	s, out := newEnv(t)
	if err := run(t, s, out, "ls", "-l", "/"); err != nil {
		t.Fatal(err)
	}
	want := "d  0 home\nd  0 etc\n- 11 readme.txt"
	if out.text() != want {
		t.Errorf("-l:\n%s\nwant:\n%s", out.text(), want)
	}
}

func TestLsAllHidden(t *testing.T) {
	s, out := newEnv(t)
	_ = run(t, s, out, "ls", "-a")
	if !strings.HasPrefix(out.text(), ". .. home") {
		t.Errorf("-a: %q", out.text())
	}
}

func TestLsRecursive(t *testing.T) {
	s, out := newEnv(t)
	if err := run(t, s, out, "ls", "-R", "/home"); err != nil {
		t.Fatal(err)
	}
	want := "/home:\nuser\n\n/home/user:\ndocuments notes.txt\n\n/home/user/documents:\ndeep.txt"
	if out.text() != want {
		t.Errorf("-R:\n%s\nwant:\n%s", out.text(), want)
	}
}

func TestLsErrors(t *testing.T) {
	s, out := newEnv(t)
	if err := run(t, s, out, "ls", "-z"); err == nil {
		t.Error("unknown flag must fail")
	}
	if err := run(t, s, out, "ls", "nope"); err == nil {
		t.Error("unknown path must fail")
	}
}

func TestHumanSize(t *testing.T) {
	cases := map[int]string{0: "0", 15: "15", 1024: "1K", 1536: "1.5K", 3 * 1024 * 1024: "3M"}
	for in, want := range cases {
		if got := humanSize(in); got != want {
			t.Errorf("humanSize(%d)=%s want %s", in, got, want)
		}
	}
}

func TestCdSharesSession(t *testing.T) {
	s, out := newEnv(t)
	if err := run(t, s, out, "cd", "/home/user"); err != nil {
		t.Fatal(err)
	}
	if err := run(t, s, out, "cd", "/home/user/documents"); err != nil {
		t.Fatal(err)
	}
	if s.PathString() != "/home/user/documents" {
		t.Errorf("path %s", s.PathString())
	}
}

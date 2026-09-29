package session

import (
	"fmt"
	"slices"

	"vfs/src/internal/vfs"
)

const (
	KindFile = "file"
	KindDir  = "directory"

	msgEntryNotFound = "file or directory %s not found"
	msgRefuseDots    = "refusing to remove '.' or '..': %s"
	msgBusyDir       = "cannot remove %s: current directory is inside it"
)

func (s *Session) Remove(target string) (string, error) {
	parent, name, err := s.locate(target)
	if err != nil {
		return "", err
	}
	if name == CurrentDir || name == ParentDir {
		return "", fmt.Errorf(msgRefuseDots, target)
	}

	dir := parent.current()
	switch {
	case removeFile(dir, name):
		return KindFile, nil
	case hasSubdir(dir, name):
		return KindDir, s.removeDir(parent, name)
	}
	return "", fmt.Errorf(msgEntryNotFound, target)
}

func (s *Session) removeDir(parent *position, name string) error {
	removed := append(slices.Clone(parent.names), name)
	if s.isInside(removed) {
		return fmt.Errorf(msgBusyDir, name)
	}

	dir := parent.current()
	dir.Subdirs = slices.DeleteFunc(dir.Subdirs, func(d vfs.Directory) bool {
		return d.Name == name
	})
	return s.syncCurrent()
}

func (s *Session) isInside(dirNames []string) bool {
	return len(dirNames) <= len(s.Path) &&
		slices.Equal(dirNames, s.Path[:len(dirNames)])
}

func (s *Session) syncCurrent() error {
	pos, err := s.start(CurrentDir)
	if err != nil {
		return err
	}
	s.Current = pos.current()
	return nil
}

func removeFile(dir *vfs.Directory, name string) bool {
	before := len(dir.Files)
	dir.Files = slices.DeleteFunc(dir.Files, func(f vfs.File) bool {
		return f.Name == name
	})
	return len(dir.Files) != before
}

func hasSubdir(dir *vfs.Directory, name string) bool {
	return findSubdir(dir, name) != nil
}

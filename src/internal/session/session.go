package session

import (
	"fmt"
	"strings"
	"vfs/src/internal/vfs"
)

type Session struct {
	Root    *vfs.Directory
	Current *vfs.Directory
	Path    []string
}

func NewSession(root *vfs.Directory) *Session {
	return &Session{
		Root:    root,
		Current: root,
		Path:    []string{},
	}
}

func (s *Session) PathString() string {
	return "/" + strings.Join(s.Path, "/")
}

func (s *Session) Cd(arg string) error {
	if arg == "" || arg == "/" {
		s.Current = s.Root
		s.Path = []string{}
		return nil
	}

	if arg == ".." {
		if len(s.Path) == 0 {
			return nil
		}
		s.Path = s.Path[:len(s.Path)-1]
		s.Current = s.findByPath(s.Path)
		return nil
	}

	if arg == "." {
		return nil
	}

	parts := strings.Split(strings.Trim(arg, "/"), "/")
	newCurrent := s.Current
	newPath := append([]string{}, s.Path...)

	for _, part := range parts {
		if part == ".." {
			if len(newPath) > 0 {
				newPath = newPath[:len(newPath)-1]
				newCurrent = s.findByPath(newPath)
			}
			continue
		}

		found := findSubdir(newCurrent, part)
		if found == nil {
			return fmt.Errorf("no such directory: %s", part)
		}
		newCurrent = found
		newPath = append(newPath, part)
	}

	s.Current = newCurrent
	s.Path = newPath
	return nil
}

func (s *Session) FindFile(path string) (*vfs.File, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	current := s.Current

	for _, part := range parts[:len(parts)-1] {
		next := findSubdir(current, part)
		if next == nil {
			return nil, fmt.Errorf("no such directory: %s", part)
		}
		current = next
	}

	fileName := parts[len(parts)-1]
	for i := range current.Files {
		if current.Files[i].Name == fileName {
			return &current.Files[i], nil
		}
	}

	return nil, fmt.Errorf("file not found: %s", fileName)
}

func findSubdir(dir *vfs.Directory, name string) *vfs.Directory {
	for i := range dir.Subdirs {
		if dir.Subdirs[i].Name == name {
			return &dir.Subdirs[i]
		}
	}
	return nil
}

func (s *Session) findByPath(path []string) *vfs.Directory {
	current := s.Root
	for _, p := range path {
		current = findSubdir(current, p)
		if current == nil {
			return s.Root
		}
	}
	return current
}

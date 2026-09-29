package vfs

import (
	"encoding/base64"
	"encoding/xml"
	"os"
)

type VFS struct {
	Name string
	Root Directory
}

type xmlVFS struct {
	XMLName xml.Name     `xml:"vfs"`
	Name    string       `xml:"name,attr"`
	Root    xmlDirectory `xml:"directory"`
}

func Load(path string) (VFS, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return VFS{}, err
	}

	var xmlVfs xmlVFS
	if err := xml.Unmarshal(data, &xmlVfs); err != nil {
		return VFS{}, err
	}

	root, err := loadDirectory(xmlVfs.Root)
	if err != nil {
		return VFS{}, err
	}

	return VFS{Name: root.Name, Root: root}, nil
}

func loadDirectory(xmlDir xmlDirectory) (Directory, error) {
	dir := Directory{
		Name: xmlDir.Name,
	}

	for _, xf := range xmlDir.Files {
		content, err := base64.StdEncoding.DecodeString(xf.Content)
		if err != nil {
			return Directory{}, err
		}
		dir.Files = append(dir.Files, File{Name: xf.Name, Content: content, Owner: "root"})
	}

	for _, xd := range xmlDir.Subdirs {
		subdir, err := loadDirectory(xd)
		if err != nil {
			return Directory{}, err
		}
		dir.Subdirs = append(dir.Subdirs, subdir)
	}

	return dir, nil
}

package vfs

type Directory struct {
	Name    string
	Files   []File
	Subdirs []Directory
}

type xmlDirectory struct {
	Name    string         `xml:"name,attr"`
	Files   []xmlFile      `xml:"file"`
	Subdirs []xmlDirectory `xml:"directory"`
}

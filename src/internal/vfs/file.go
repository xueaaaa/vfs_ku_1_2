package vfs

type File struct {
	Name    string
	Content []byte
	Owner   string
}

type xmlFile struct {
	Name    string `xml:"name,attr"`
	Content string `xml:"content,attr"`
	Owner   string `xml:"-"`
}

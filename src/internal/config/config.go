package config

import "flag"

type Config struct {
	VfsPath    string
	ScriptPath string
}

func New() Config {
	vfsPath := flag.String("vfs", "", "path to vfs")
	scriptPath := flag.String("script", "", "path to initial script")
	flag.Parse()

	return Config{
		VfsPath:    *vfsPath,
		ScriptPath: *scriptPath,
	}
}

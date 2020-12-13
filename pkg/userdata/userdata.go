package userdata

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

const (
	typeFile = "file"
	typeDir  = "dir"

	defaultFileMode = "0644"
	defaultDirMode  = "0755"
	defaultType     = typeFile

	parenmode = os.ModeDir | 0755
)

// Config represents the userdata configuraion.
type Config struct {
	Files []File `yaml:"files"`
}

// File represents a single file system entry to be added.
type File struct {
	mode os.FileMode
	Mode string `yaml:"mode"`
	Type string `yaml:"type"`
	Path string `yaml:"path"`
	Text string `yaml:"text"`
	Link string `yaml:"link"`
}

// Apply apllies changes described in the userdata config provided.
func Apply(prefix string, c *Config) error {
	if err := normalize(c.Files); err != nil {
		return err
	}

	return apply(prefix, c.Files)
}

var types = map[string]bool{typeDir: true, typeFile: true}

func ferrorf(i int, emsg string, err error) error {
	if err != nil {
		return fmt.Errorf("file[%d]: %s: %v", i, emsg, err)
	}
	return fmt.Errorf("file[%d]: %s", i, emsg)
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}

func normalize(files []File) error {
	for i := 0; i < len(files); i++ {
		if files[i].Type == "" {
			files[i].Type = defaultType
		}
		if _, ok := types[files[i].Type]; !ok {
			return ferrorf(i, "unexpected type", nil)
		}
		if files[i].Mode == "" {
			if files[i].Type == typeDir {
				files[i].Mode = defaultDirMode
			} else {
				files[i].Mode = defaultFileMode
			}
		}
		mode, err := strconv.ParseUint(files[i].Mode, 8, 32)
		if err != nil {
			return ferrorf(i, "invalid mode value", err)
		}
		files[i].mode = os.FileMode(uint32(mode))
	}
	return nil
}

func apply(prefix string, files []File) error {
	for i, f := range files {
		if err := os.MkdirAll(path.Dir(f.Path), parenmode); err != nil {
			return ferrorf(i, "unable to create parent directory", err)
		}
		if f.Type == typeFile {
			if err := ioutil.WriteFile(f.Path, []byte(f.Text), f.mode); err != nil {
				return ferrorf(i, "unable to write file", err)
			}
		} else {
			if err := os.MkdirAll(f.Path, os.ModeDir|f.mode); err != nil {
				return ferrorf(i, "unable to create directory", err)
			}
		}
		if f.Link != "" {
			if err := os.MkdirAll(path.Dir(f.Path), parenmode); err != nil {
				return ferrorf(i, "unable to create link parent directory", err)
			}
			if err := os.RemoveAll(f.Link); err != nil {
				return ferrorf(i, "cannot remove link", err)
			}
			if err := os.Symlink(filepath.Join(prefix, f.Path), f.Link); err != nil {
				return ferrorf(i, "unable to link", err)
			}
		}
	}
	return nil
}

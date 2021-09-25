package vdir

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

// Collection represents a Vdir collection
type Collection struct {
	Name  string
	Color string
	Path  string
}

const (
	MetaDisplayName = "displayname" // MetaDisplayName is a filename vdir uses for collection name
	MetaColor       = "color"       // MetaColor is a filename vdir uses for collection color
)

func (c *Collection) String() string {
	return c.Name
}

// Init initializes a Collection with a path, name and color parsed from path
func (c *Collection) Init(path string) error {
	f, err := os.Stat(path)
	if errors.Is(err, fs.ErrNotExist) {
		return errors.New("Specified path does not exist")
	}
	if !f.IsDir() {
		return errors.New("Specified path is not a directory")
	}

	c.Path = path
	c.Name = filepath.Base(path)

	// parse dir for metadata
	err = filepath.WalkDir(
		path,
		func(pp string, dd fs.DirEntry, err error) error {
			if dd.Name() == MetaDisplayName {
				name, err := os.ReadFile(pp)
				if err != nil {
					return err
				}
				c.Name = string(name)
			}
			if dd.Name() == MetaColor {
				color, err := os.ReadFile(pp)
				if err != nil {
					return err
				}
				c.Color = string(color)
			}
			return nil
		},
	)
	if err != nil {
		return err
	}
	return nil
}

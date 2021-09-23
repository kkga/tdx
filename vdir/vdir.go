package vdir

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/emersion/go-ical"
)

// VdirRoot represents the topmost vdir root folder
type VdirRoot struct {
	Path string
}

// Collections represents a map all collections and their items
type Collections map[*Collection][]*Item

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

// NewVdirRoot initializes a VdirRoot checking that the directory exists
func NewVdirRoot(path string) (*VdirRoot, error) {
	f, err := os.Stat(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, errors.New("Specified vdir path does not exist")
	}
	if !f.IsDir() {
		return nil, errors.New("Specified vdir path is not a directory")
	}
	return &VdirRoot{path}, nil
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

func (c Collection) String() string {
	return c.Name
}

// Collections returns a map of all collections and items in vdir, items have unique id values
func (v VdirRoot) Collections() (collections Collections, err error) {
	collections = make(Collections)
	id := 0

	isIcal := func(path string, de fs.DirEntry) bool {
		return !de.IsDir() && strings.TrimPrefix(filepath.Ext(path), ".") == ical.Extension
	}

	hasIcalFiles := func(path string) bool {
		files, _ := os.ReadDir(path)
		for _, f := range files {
			if filepath.Ext(f.Name()) == fmt.Sprintf(".%s", ical.Extension) {
				return true
			}
		}
		return false
	}

	// parse dir for vdir collections (folders)
	err = filepath.WalkDir(
		v.Path,
		func(p string, d fs.DirEntry, err error) error {
			if d.IsDir() && hasIcalFiles(p) {
				c := &Collection{}
				if err := c.Init(p); err != nil {
					return err
				}
				// parse collection folder for ical files
				err = filepath.WalkDir(c.Path, func(pp string, dd fs.DirEntry, err error) error {
					if isIcal(pp, dd) {
						item := &Item{}
						item.Init(p)
						if err := item.Init(pp); err != nil {
							return err
						}
						if item.Ical != nil {
							id++
							item.Id = id
							collections[c] = append(collections[c], item)
						}
					}
					return nil
				})
				if err != nil {
					return err
				}
			}
			return nil
		},
	)
	if err != nil {
		return
	}
	return collections, nil
}

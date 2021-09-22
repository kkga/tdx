package vdir

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/emersion/go-ical"
)

func NewVdirRoot(path string) (*VdirRoot, error) {
	f, err := os.Stat(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, errors.New("Specified vdir path does not exist.")
	}
	if !f.IsDir() {
		return nil, errors.New("Specified vdir path is not a directory.")
	}
	return &VdirRoot{path}, nil
}

// VdirRoot represents the topmost vdir root folder
type VdirRoot struct {
	Path string
}

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

// Collections returns a slice of all vdir collections in root path recursively
func (v VdirRoot) Collections() (collections []*Collection, err error) {
	err = filepath.WalkDir(v.Path, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() && p != v.Path {
			var c = &Collection{}
			c.Path = p
			c.Name = d.Name()

			err = filepath.WalkDir(p, func(pp string, dd fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
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
			})
			if err != nil {
				return err
			}
			collections = append(collections, c)
		}
		return nil
	})
	if err != nil {
		return
	}
	return collections, nil
}

// Items returns a slice of decoded iCalendar items in collection
func (c Collection) Items() (items []*ical.Calendar, err error) {

	isIcal := func(path string, de fs.DirEntry) bool {
		return !de.IsDir() && strings.TrimPrefix(filepath.Ext(path), ".") == ical.Extension
	}

	err = filepath.WalkDir(c.Path, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if isIcal(p, d) {
			file, err := os.Open(p)
			if err != nil {
				return err
			}
			defer file.Close()

			dec := ical.NewDecoder(file)

			for {
				item, err := dec.Decode()
				if err == io.EOF {
					break
				} else if err != nil {
					return err
				}

				// filter only items that contain vtodo
				for _, comp := range item.Children {
					if comp.Name == ical.CompToDo {
						items = append(items, item)
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		return
	}
	return items, nil
}

func (c Collection) String() string {
	return c.Name
}

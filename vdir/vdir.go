package vdir

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/emersion/go-ical"
)

// Collection represents a Vdir collection
type Collection struct {
	Name  string
	Color string
	Path  string
}

const (
	MetaDisplayName = "displayname"
	MetaColor       = "color"
)

// Collections returns a list of all vdir collections in specified root path recursively
func Collections(root string) (collections []Collection, err error) {
	err = filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() && p != root {
			var c = &Collection{}
			c.Path = p

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
			collections = append(collections, *c)
		}
		return nil
	})
	if err != nil {
		return
	}
	return collections, nil
}

// Items returns a slice of all iCalendar items in collection
func (c Collection) Items() (items []ical.Calendar, err error) {

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
				items = append(items, *item)
			}
		}
		return nil
	})
	if err != nil {
		return
	}
	return items, nil
}

// The vdir package implements a simple library to interact with vdir
// filesystem storage: https://vdirsyncer.pimutils.org/en/stable/vdir.html
package vdir

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/emersion/go-ical"
)

// Vdir is a map of all collections and items
type Vdir map[*Collection][]*Item

// Init initializes the map with collections and items in path, items have unique IDs
func (v *Vdir) Init(path string) error {
	f, err := os.Stat(path)
	if errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("Vdir path does not exist: %q", path)
	}
	if !f.IsDir() {
		return fmt.Errorf("Vdir path is not a directory: %q", path)
	}

	var itemId int

	walkFunc := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && hasIcalFiles(path) {
			c := &Collection{}
			if err := c.Init(path); err != nil {
				return err
			}
			// parse collection folder for ical files
			err = filepath.WalkDir(
				c.Path,
				func(pp string, dd fs.DirEntry, err error) error {
					if isIcal(pp, dd) {
						item := new(Item)
						if err := item.Init(pp); err != nil {
							log.Println(err)
						}
						if item.Ical != nil {
							itemId++
							item.Id = itemId
							(*v)[c] = append((*v)[c], item)
						}
					}
					return nil
				},
			)
			if err != nil {
				return err
			}
		}
		return nil
	}

	err = filepath.WalkDir(path, walkFunc)
	if err != nil {
		return err
	}

	return nil
}

// ItemById finds and returns an item for specified id
func (v *Vdir) ItemById(id int) (*Item, error) {
	for _, items := range *v {
		for _, item := range items {
			if item.Id == id {
				return item, nil
			}
		}
	}
	return nil, fmt.Errorf("Item not found: %d", id)
}

// ItemByPath finds and returns an item for specified path
func (v *Vdir) ItemByPath(path string) (*Item, error) {
	for _, items := range *v {
		for _, item := range items {
			if item.Path == path {
				return item, nil
			}
		}
	}
	return nil, fmt.Errorf("Item not found: %q", path)
}

// Tags returns a slice of all tags found in todos inside vdir
func (v *Vdir) Tags() (tags []Tag, err error) {
	containsTag := func(tags []Tag, tag Tag) bool {
		for _, t := range tags {
			if tag == t {
				return true
			}
		}
		return false
	}

	for _, items := range *v {
		for _, item := range items {
			tt, err := item.Tags()
			if err != nil {
				return tags, err
			}
			for _, tag := range tt {
				if !containsTag(tags, tag) {
					tags = append(tags, tag)
				}
			}
		}
	}

	return
}

// isIcal reports whether path is a file that has an ical extension
func isIcal(path string, de fs.DirEntry) bool {
	return !de.IsDir() && filepath.Ext(path) == fmt.Sprintf(".%s", ical.Extension)
}

// hasIcalFiles reports whether path contains ical files
func hasIcalFiles(path string) bool {
	files, _ := os.ReadDir(path)
	for _, f := range files {
		if filepath.Ext(f.Name()) == fmt.Sprintf(".%s", ical.Extension) {
			return true
		}
	}
	return false
}

package vdir

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/emersion/go-ical"
)

// VdirRoot is topmost vdir root folder
type VdirRoot struct {
	Path string
}

// VdirMap is a map all collections and their items
type VdirMap map[*Collection][]*Item

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

// ItemById finds and returns an item for specified id
func (c VdirMap) ItemById(id int) (*Item, error) {
	for _, items := range c {
		for _, item := range items {
			if item.Id == id {
				return item, nil
			}
		}
	}
	return nil, fmt.Errorf("Item not found: %d", id)
}

// Collections returns a map of all collections and items in vdir, items have unique id values
func (v VdirRoot) InitMap() (vdirMap VdirMap, err error) {
	vdirMap = make(VdirMap)

	isIcal := func(path string, de fs.DirEntry) bool {
		return !de.IsDir() && filepath.Ext(path) == fmt.Sprintf(".%s", ical.Extension)
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

	var itemId int
	walkFunc := func(path string, d fs.DirEntry, err error) error {
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
							return err
						}
						if item.Ical != nil {
							itemId++
							item.Id = itemId
							vdirMap[c] = append(vdirMap[c], item)
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

	// parse dir for vdir collections (folders)
	err = filepath.WalkDir(v.Path, walkFunc)
	if err != nil {
		return
	}

	return vdirMap, nil
}

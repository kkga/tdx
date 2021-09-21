package vdir

import (
	"io"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/emersion/go-ical"
)

func decodeFiles(files []fs.FileInfo) []ical.Component {
	components := []ical.Component{}

	filterComponents := func(cal *ical.Calendar, filter string) []ical.Component {
		l := make([]ical.Component, 0, len(cal.Children))
		for _, child := range cal.Children {
			if child.Name == filter {
				l = append(l, *child)
			}
		}
		return l
	}

	for _, f := range files {
		if strings.TrimPrefix(filepath.Ext(f.Name()), ".") != ical.Extension {
			continue
		}

		path := path.Join(calDir, f.Name())
		file, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		dec := ical.NewDecoder(file)

		for {
			cal, err := dec.Decode()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}

			components = append(components, filterComponents(cal, ical.CompToDo)...)
		}
	}

	return components
}

package main

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/emersion/go-ical"
	"github.com/kkga/ctdo/cmd"
)

var calDir = "/home/kkga/.local/share/calendars/tasks/"
var version = "dev"

func main() {
	// log.SetFlags(0)

	if len(os.Args) > 1 && os.Args[1] == "-v" {
		fmt.Printf("kks %s\n", version)
		os.Exit(0)
	}

	if err := cmd.Root(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
	// if len(os.Args) > 1 && os.Args[1] == "new" {
	// 	todo := &todo.ToDo{}
	// 	todoBuf, _ := todo.Encode()

	// 	f, _ := os.Create(fmt.Sprintf("%s/ctdo-testing-%d.ics", calDir, time.Now().UnixNano()))
	// 	defer f.Close()
	// 	w := bufio.NewWriter(f)
	// 	_, _ = w.Write(todoBuf.Bytes())
	// 	w.Flush()
	// } else {
	// 	files, err := ioutil.ReadDir(calDir)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	todos := decodeFiles(files)
	// 	list := todo.NewList()
	// 	list.Init("new list", todos)
	// 	fmt.Println(list.String())
	// }
}

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

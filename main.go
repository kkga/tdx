package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/emersion/go-ical"
)

var calDir = "/home/kkga/.local/share/vdirsyncer/migadu-cal/tasks/"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "new" {
		encode()
	} else {
		list := list()
		printTodos(list)
	}
}

func encode() {
	event := ical.NewEvent()
	event.Props.SetText(ical.PropUID, "uid@example.org")
	event.Props.SetDateTime(ical.PropDateTimeStamp, time.Now())
	event.Props.SetText(ical.PropSummary, "My awesome event")
	event.Props.SetDateTime(ical.PropDateTimeStart, time.Now().Add(24*time.Hour))

	cal := ical.NewCalendar()
	cal.Props.SetText(ical.PropVersion, "2.0")
	cal.Props.SetText(ical.PropProductID, "-//xyz Corp//NONSGML PDA Calendar Version 1.0//EN")
	cal.Children = append(cal.Children, event.Component)

	var buf bytes.Buffer
	if err := ical.NewEncoder(&buf).Encode(cal); err != nil {
		log.Fatal(err)
	}

	log.Print(buf.String())
}

func list() []ToDo {
	files, err := ioutil.ReadDir(calDir)
	if err != nil {
		log.Fatal(err)
	}

	list := make([]ToDo, 0, len(files))

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

			t := todos(cal)

			for _, todo := range t {
				t := NewToDo()
				err := t.ParseComponent(todo)
				if err != nil {
					log.Fatal(err)
				}
				list = append(list, *t)
			}
		}
	}

	return list
}

func todos(cal *ical.Calendar) []ical.Component {
	l := make([]ical.Component, 0, len(cal.Children))
	for _, child := range cal.Children {
		if child.Name == ical.CompToDo {
			l = append(l, *child)
		}
	}
	return l
}

func printTodos(todos []ToDo) {
	for _, t := range todos {
		if t.Status != ToDoStatusCompleted {
			fmt.Println(t.String())
		}
	}
}

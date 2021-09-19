package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/emersion/go-ical"
)

var calDir = "/home/kkga/.local/share/vdirsyncer/migadu-cal/tasks/"

func main() {
	list := []ToDo{}

	files, err := ioutil.ReadDir(calDir)
	if err != nil {
		log.Fatal(err)
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

			t := todos(cal)

			for _, todo := range t {
				t, err := NewToDo(todo)
				if err != nil {
					log.Fatal(err)
				}
				list = append(list, *t)
			}
		}
	}

	printTodos(list)
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

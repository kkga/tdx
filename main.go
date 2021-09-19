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
	todos := []ToDo{}

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

			for _, comp := range cal.Children {
				if comp.Name == ical.CompToDo {
					t := ToDo{}

					// fmt.Println("-----------")

					for p := range comp.Props {

						// fmt.Println(p)
						switch p {
						case ical.PropStatus:
							s, err := comp.Props.Get(ical.PropStatus).Text()
							if err != nil {
								log.Fatal(err)
							}
							t.Status = ToDoStatus(s)
						case ical.PropSummary:
							s, err := comp.Props.Get(ical.PropSummary).Text()
							if err != nil {
								log.Fatal(err)
							}
							t.Summary = s
						case ical.PropDescription:
							s, err := comp.Props.Get(ical.PropDescription).Text()
							if err != nil {
								log.Fatal(err)
							}
							t.Description = s
						case ical.PropDue:
							time, err := comp.Props.Get(ical.PropDue).DateTime(t.Due.Location())
							if err != nil {
								log.Fatal(err)
							}
							t.Due = time
						case ical.PropPriority:
							prio, err := comp.Props.Get(ical.PropPriority).Int()
							if err != nil {
								log.Fatal(err)
							}
							t.Priority = prio
						}
					}

					todos = append(todos, t)
				}
			}
		}
	}

	printTodos(todos)
}

func printTodos(todos []ToDo) {
	for _, t := range todos {
		if t.Status != ToDoStatusCompleted {
			fmt.Println(t.String())
		}
	}
}

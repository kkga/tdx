package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/emersion/go-ical"
)

var calDir = "/home/kkga/.local/share/vdirsyncer/migadu-cal/tasks/"

func main() {
	todos := []Todo{}

	files, _ := ioutil.ReadDir(calDir)

	for _, f := range files {
		file, err := os.Open(fmt.Sprintf("%s/%s", calDir, f.Name()))
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
					t := Todo{}

					fmt.Println("--")
					for p := range comp.Props {
						fmt.Println(p)
						switch p {
						case ical.PropStatus:
							s, err := comp.Props.Get(ical.PropStatus).Text()
							if err != nil {
								log.Fatal(err)
							}
							t.Status = TodoStatus(s)
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
						}
					}

					todos = append(todos, t)
				}
			}
		}
	}

	for _, t := range todos {
		fmt.Println(t.String())
	}
}

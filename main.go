package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/emersion/go-ical"
)

var calPath = "/home/kkga/.local/share/vdirsyncer/migadu-cal/tasks/0cd4d8f7-4acc-40cc-bc51-752fded1fa7d.ics"

func main() {
	// var r io.Reader
	f, _ := os.Open(calPath)
	defer f.Close()

	dec := ical.NewDecoder(f)

	for {
		cal, err := dec.Decode()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		for _, c := range cal.Children {
			// fmt.Println(c.Name)
			for _, p := range c.Props {
				name := p[0].Name
				val := p[0].Value
				switch p[0].Name {
				case ical.PropStatus:
					s := p[0].Value
					if s == "NEEDS-ACTION" {
						val = "**not done**"
					} else {
						val = "**done**"
					}
				}
				fmt.Println(name, val)
				// fmt.Printf("%+v\n", p[0].Text())
			}
		}

		for _, event := range cal.Events() {
			summary, err := event.Props.Text(ical.PropSummary)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Found event: %v", summary)
		}
	}
}

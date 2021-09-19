package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/emersion/go-ical"
)

// var calPath = "/home/kkga/.local/share/vdirsyncer/migadu-cal/tasks/0cd4d8f7-4acc-40cc-bc51-752fded1fa7d.ics"
var calPath = "/home/kkga/.local/share/vdirsyncer/migadu-cal/tasks/0FE5B89B-57A8-4069-93A6-0A412334A33F.ics"

var calDir = "/home/kkga/.local/share/vdirsyncer/migadu-cal/tasks"

func main() {

	// var r io.Reader
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

			for _, c := range cal.Children {
				summary := c.Props.Get(ical.PropSummary)
				desc := c.Props.Get(ical.PropDescription)
				fmt.Println(summary, desc)
				// fmt.Println(c.Name)
				// for _, p := range c.Props {
				// 	name := p[0].Name
				// 	val := p[0].Value
				// 	switch p[0].Name {
				// 	case ical.PropStatus:
				// 		s := p[0].Value
				// 		if s == "NEEDS-ACTION" {
				// 			val = "**not done**"
				// 		} else {
				// 			val = "**done**"
				// 		}
				// 	}
				// 	fmt.Println(name, val)
				// fmt.Printf("%+v\n", p[0].Text())
				// }
			}
			fmt.Println("-------------------------")

			for _, event := range cal.Events() {
				summary, err := event.Props.Text(ical.PropSummary)
				if err != nil {
					log.Fatal(err)
				}
				log.Printf("Found event: %v", summary)
			}
		}

	}

}

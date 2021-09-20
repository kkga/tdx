package cmd

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/emersion/go-ical"
)

type Runner interface {
	Init([]string) error
	Run() error
	Name() string
	Alias() []string
}

type Cmd struct {
	fs        *flag.FlagSet
	alias     []string
	shortDesc string
	usageLine string
	cals      []ical.Calendar
}

var calDir = "/home/kkga/.local/share/calendars/tasks/"

func (c *Cmd) Run() error      { return nil }
func (c *Cmd) Name() string    { return c.fs.Name() }
func (c *Cmd) Alias() []string { return c.alias }

func (c *Cmd) Init(args []string) error {
	c.fs.Usage = c.usage

	if err := c.fs.Parse(args); err != nil {
		return err
	}

	return nil
}

func (c *Cmd) usage() {
	fmt.Println(c.shortDesc)
	fmt.Println()

	fmt.Println("USAGE")
	fmt.Printf("  kks %s %s\n\n", c.fs.Name(), c.usageLine)

	if strings.Contains(c.usageLine, "[options]") {
		fmt.Println("OPTIONS")
		c.fs.PrintDefaults()
	}
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

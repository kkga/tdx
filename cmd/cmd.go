package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/emersion/go-ical"
	"github.com/kkga/tdx/vdir"
)

type Runner interface {
	Init([]string) error
	Run() error
	Name() string
	Alias() []string
}

type Cmd struct {
	fs         *flag.FlagSet
	alias      []string
	shortDesc  string
	usageLine  string
	cals       []ical.Calendar
	root       *vdir.VdirRoot
	collection *vdir.Collection

	listFlag string
	listReq  bool
}

var dir = "/home/kkga/.local/share/calendars/migadu/"

func (c *Cmd) Run() error      { return nil }
func (c *Cmd) Name() string    { return c.fs.Name() }
func (c *Cmd) Alias() []string { return c.alias }

func (c *Cmd) Init(args []string) error {
	env := struct{ list string }{list: os.Getenv("TDX_DEFAULT_LIST")}

	if env.list != "" {
		c.listFlag = env.list
	}

	c.fs.Usage = c.usage

	if err := c.fs.Parse(args); err != nil {
		return err
	}

	root, err := vdir.NewVdirRoot(dir)
	if err != nil {
		return err
	}
	c.root = root

	if c.listReq && c.listFlag == "" {
		return errors.New("Specify a list with '-l' or set default list with 'TDX_DEFAULT_LIST'")
	} else {
		collections, err := root.Collections()
		if err != nil {
			return err
		}

		names := []string{}

		for col := range collections {
			names = append(names, col.Name)
			if col.Name == c.listFlag {
				c.collection = col
			}
		}
		if c.collection == nil {
			return fmt.Errorf("List does not exist: %s\nAvailable lists: %s", c.listFlag, strings.Join(names, ", "))
		}
	}

	return nil
}

func (c *Cmd) usage() {
	fmt.Println(c.shortDesc)
	fmt.Println()

	fmt.Println("USAGE")
	fmt.Printf("  ctodo %s %s\n\n", c.fs.Name(), c.usageLine)

	if strings.Contains(c.usageLine, "[options]") {
		fmt.Println("OPTIONS")
		c.fs.PrintDefaults()
	}
}

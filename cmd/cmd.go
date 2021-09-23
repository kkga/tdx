package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/kkga/tdx/vdir"
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
	// cals       []ical.Calendar

	root           *vdir.VdirRoot
	allCollections map[*vdir.Collection][]*vdir.Item
	collection     *vdir.Collection

	list    string
	listReq bool
}

var dir = "/home/kkga/.local/share/calendars/migadu/"

func (c *Cmd) Run() error      { return nil }
func (c *Cmd) Name() string    { return c.fs.Name() }
func (c *Cmd) Alias() []string { return c.alias }

func (c *Cmd) Init(args []string) error {
	env := struct{ list string }{list: os.Getenv("TDX_DEFAULT_LIST")}

	if env.list != "" {
		c.list = env.list
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

	collections, err := root.Collections()
	if err != nil {
		return err
	}
	c.allCollections = collections

	if c.listReq && c.list == "" {
		return errors.New("Specify a list with '-l' or set default list with 'TDX_DEFAULT_LIST'")
	} else if c.list != "" {

		names := []string{}
		for col := range collections {
			names = append(names, col.Name)
			if col.Name == c.list {
				c.collection = col
			}
		}
		if c.collection == nil {
			return fmt.Errorf("List does not exist: %s\nAvailable lists: %s", c.list, strings.Join(names, ", "))
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

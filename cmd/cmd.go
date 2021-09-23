package cmd

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/kelseyhightower/envconfig"
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

	root           *vdir.VdirRoot
	allCollections vdir.Collections
	collection     *vdir.Collection

	listFlag     string
	listRequired bool
}

type Config struct {
	Path          string `required:"true"`
	DefaultList   string `split_words:"true"`
	DefaultStatus string `split_words:"true" default:"NEEDS-ACTION"`
	DefaultSort   string `split_words:"true" default:"PRIORITY"`
	DefaultDue    int    `split_words:"true" default:"48"`
	Color         bool   `default:"true"`
}

func (c *Cmd) Run() error      { return nil }
func (c *Cmd) Name() string    { return c.fs.Name() }
func (c *Cmd) Alias() []string { return c.alias }

func (c *Cmd) Init(args []string) error {
	var conf Config
	err := envconfig.Process("TDX", &conf)
	if err != nil {
		return err
	}

	c.listFlag = conf.DefaultList
	c.fs.Usage = c.usage

	if err := c.fs.Parse(args); err != nil {
		return err
	}

	root, err := vdir.NewVdirRoot(conf.Path)
	if err != nil {
		return err
	}
	c.root = root

	collections, err := root.Collections()
	if err != nil {
		return err
	}
	c.allCollections = collections

	if c.listRequired && c.listFlag == "" {
		return errors.New("Specify a list with '-l' or set default list with 'TDX_DEFAULT_LIST'")
	} else if c.listFlag != "" {

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
	fmt.Printf("  tdx %s %s\n\n", c.fs.Name(), c.usageLine)

	if strings.Contains(c.usageLine, "[options]") {
		fmt.Println("OPTIONS")
		c.fs.PrintDefaults()
	}
}
